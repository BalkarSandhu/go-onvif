package api

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-onvif/json_apis/utils"
	"go-onvif/onvif"
	wsdiscovery "go-onvif/ws-discovery"

	"github.com/beevik/etree"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

const (
	scanTimeout = 10 * time.Second
)

type DeviceScanResult struct {
	IP              string `json:"ip"`
	Alive           bool   `json:"alive"`
	Manufacturer    string `json:"manufacturer,omitempty"`
	Model           string `json:"model,omitempty"`
	FirmwareVersion string `json:"firmware_version,omitempty"`
	SerialNumber    string `json:"serial_number,omitempty"`
	HardwareID      string `json:"hardware_id,omitempty"`
	Error           string `json:"error,omitempty"`
}

type ScanRequest struct {
	StartIP  string `form:"start" binding:"required"`
	EndIP    string `form:"end" binding:"required"`
	Username string `form:"username"`
	Password string `form:"password"`
	Timeout  int    `form:"timeout"`
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables directly...")
	}

	getEnv := func(key, fallback string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return fallback
	}

	toInt := func(val string, fallback int) int {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		return fallback
	}

	return Config{
		Port:           getEnv("PORT", "8084"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RateLimitReqs:  toInt(getEnv("RATE_LIMIT_REQS", "10"), 10),
		RateLimitBurst: toInt(getEnv("RATE_LIMIT_BURST", "20"), 20),
	}
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
	gin.SetMode(gin.DebugMode)
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

	//Scan IP range
	s.router.GET("/scan", s.handleDeviceScan)
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
	response, err := s.CallServiceMethod(serviceName, methodName, acceptedData, dev)
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
func (s *APIServer) CallServiceMethod(serviceName, methodName string, data []byte, dev *onvif.Device) (interface{}, error) {
	switch strings.ToLower(serviceName) {
	case "device":
		return utils.CallDeviceMethod(methodName, dev, data)
	case "ptz":
		return utils.CallPTZMethod(methodName, dev, data)
	case "media":
		return utils.CallPTZMethod(methodName, dev, data)
	default:
		return nil, errors.New("unknown service: " + serviceName)
	}
}

// handleDiscovery handles device discovery requests
func (s *APIServer) handleDiscovery(c *gin.Context) {
	interfaceName := c.Query("interface")

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

func (s *APIServer) handleDeviceScan(c *gin.Context) {
	var req utils.ScanRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		s.logger.Error().Err(err).Msg("Invalid scan request parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Use default timeout if not set
	if req.Timeout == 0 {
		req.Timeout = int(scanTimeout.Seconds()) // assumes scanTimeout is a predefined `time.Duration`
	}
	scanTimeoutDuration := time.Duration(req.Timeout) * time.Second

	// Generate list of IPs to scan
	ips, err := utils.GenerateIPRange(req.StartIP, req.EndIP)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate IP range")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid IP range: " + err.Error(),
		})
		return
	}

	if len(ips) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "IP range too large (max 1000 addresses)",
		})
		return
	}

	// Get body data if required for auth headers, etc.
	acceptedData, err := c.GetRawData()
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to read raw request body")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process request data",
		})
		return
	}

	s.logger.Info().
		Str("start_ip", req.StartIP).
		Str("end_ip", req.EndIP).
		Int("ip_count", len(ips)).
		Bool("with_auth", req.Username != "").
		Dur("timeout", scanTimeoutDuration).
		Msg("Starting ONVIF device scan")

	// Perform scan
	var scanResult utils.DeviceScanResult
	results := scanResult.PerformDeviceScan(ips, req, acceptedData, scanTimeoutDuration)
	aliveCount := utils.CountAliveDevices(results)

	s.logger.Info().
		Int("total_scanned", len(results)).
		Int("alive_devices", aliveCount).
		Msg("ONVIF device scan completed")

	// Return results
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"summary": gin.H{
			"total_scanned": len(results),
			"alive_devices": aliveCount,
		},
	})
}

// Run starts the API server
func (s *APIServer) Run() error {
	s.logger.Info().Str("port", s.config.Port).Msg("Starting ONVIF API server")
	return s.router.Run(":" + s.config.Port)
}
