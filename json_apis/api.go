package api

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/BalkarSandhu/go-onvif/json_apis/utils"
	"github.com/BalkarSandhu/go-onvif/onvif"
	wsdiscovery "github.com/BalkarSandhu/go-onvif/ws-discovery"

	"github.com/beevik/etree"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// Config holds application configuration
type Config struct {
	Port           string
	LogLevel       string
	RateLimitReqs  int
	RateLimitBurst int
}

// APIServer is the main server structure
type APIServer struct {
	router      *gin.Engine
	logger      zerolog.Logger
	deviceCache *utils.DeviceCache
	limiter     *rate.Limiter
	config      Config
}

// NewAPIServer creates a new API server
func NewAPIServer(config Config) *APIServer {
	// Set up zerolog
	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logContext := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp()

	logger := logContext.Logger().Level(logLevel)

	// Configure gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Set up CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "username", "password", "xaddr"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create rate limiter
	limiter := rate.NewLimiter(rate.Limit(config.RateLimitReqs), config.RateLimitBurst)

	return &APIServer{
		router:      router,
		logger:      logger,
		deviceCache: utils.NewDeviceCache(10 * time.Minute),
		limiter:     limiter,
		config:      config,
	}
}

// rateLimitMiddleware provides basic rate limiting
func (s *APIServer) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// loggerMiddleware provides request logging
func (s *APIServer) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after request is complete
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		s.logger.Info().
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Msg("Request processed")
	}
}

// SetupRoutes configures all API routes
func (s *APIServer) SetupRoutes() {
	s.router.Use(s.loggerMiddleware())
	s.router.Use(s.rateLimitMiddleware())

	// Service endpoints
	s.router.POST("/:service/:method", s.handleServiceMethod)

	// Discovery endpoint
	s.router.GET("/discovery", s.handleDiscovery)
}

// handleServiceMethod processes all ONVIF service method calls
func (s *APIServer) handleServiceMethod(c *gin.Context) {
	serviceName := c.Param("service")
	methodName := c.Param("method")
	username := c.GetHeader("username")
	password := c.GetHeader("password")
	xaddr := c.GetHeader("xaddr")

	if xaddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing xaddr header",
		})
		return
	}

	acceptedData, err := c.GetRawData()
	if err != nil {
		s.logger.Debug().Err(err).Msg("Failed to get raw data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Get device from cache or create new connection
	dev, err := s.deviceCache.GetDevice(xaddr, username, password)
	if err != nil {
		s.logger.Error().Err(err).
			Str("xaddr", xaddr).
			Msg("Failed to connect to device")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to connect to device: " + err.Error(),
		})
		return
	}

	// Call the appropriate service method
	response, err := s.callServiceMethod(serviceName, methodName, acceptedData, dev)
	if err != nil {
		s.logger.Error().Err(err).
			Str("service", serviceName).
			Str("method", methodName).
			Msg("Method call failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// callServiceMethod routes the call to the appropriate service
func (s *APIServer) callServiceMethod(serviceName, methodName string, data []byte, dev *onvif.Device) (interface{}, error) {
	switch strings.ToLower(serviceName) {
	case "device":
		return callDeviceMethod(methodName, dev, data)
	case "ptz":
		return callPTZMethod(methodName, dev, data)
	case "media":
		return callMediaMethod(methodName, dev, data)
	default:
		return nil, errors.New("unknown service: " + serviceName)
	}
}

// handleDiscovery handles device discovery requests
func (s *APIServer) handleDiscovery(c *gin.Context) {
	interfaceName := c.GetHeader("interface")

	devices, err := wsdiscovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		s.logger.Error().Err(err).Msg("Discovery failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Discovery failed: " + err.Error(),
		})
		return
	}

	discoveredDevices := []map[string]string{}

	for _, deviceXML := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(deviceXML); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to parse device XML")
			continue
		}

		endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
		scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")

		if len(endpoints) == 0 {
			continue
		}

		// Get the device URL
		xaddrFull := strings.Split(endpoints[0].Text(), " ")[0]
		xaddr := strings.Split(xaddrFull, "/")[2]

		// Skip if we've already found this device
		alreadyFound := false
		for _, device := range discoveredDevices {
			if device["url"] == xaddr {
				alreadyFound = true
				break
			}
		}
		if alreadyFound {
			continue
		}

		// Extract device name from scopes
		deviceName := "Unknown"
		if len(scopes) > 0 {
			re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/([A-Za-z0-9-]+)`)
			matches := re.FindStringSubmatch(scopes[0].Text())
			if len(matches) > 1 {
				deviceName = matches[1]
			}
		}

		discoveredDevices = append(discoveredDevices, map[string]string{
			"url":  xaddr,
			"name": deviceName,
		})
	}

	c.JSON(http.StatusOK, discoveredDevices)
}

// Run starts the API server
func (s *APIServer) Run() error {
	s.logger.Info().Str("port", s.config.Port).Msg("Starting ONVIF API server")
	return s.router.Run(":" + s.config.Port)
}
