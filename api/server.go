package api

import (
	"os"
	"time"

	"go-onvif/api/cache"
	"go-onvif/api/config"
	"go-onvif/api/service"
	"go-onvif/api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type APIServer struct {
	router        *gin.Engine
	service				service.Service
	logger        zerolog.Logger
	deviceCache   *cache.DeviceCache
	deviceScanner utils.DeviceScanner
	limiter       *rate.Limiter
	config        config.Config
}

func NewAPIServer(cfg config.Config) *APIServer {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger().Level(logLevel)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.SetTrustedProxies([]string{})
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "username", "password", "xaddr"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	limiter := rate.NewLimiter(rate.Limit(cfg.RateLimitReqs), cfg.RateLimitBurst)

	cache.GlobalDeviceCache = cache.NewDeviceCache(10 * time.Minute)

	return &APIServer{
		router:      router,
		logger:      logger,
		deviceCache: cache.GlobalDeviceCache,
		limiter:     limiter,
		config:      cfg,
	}
}

// Run starts the API server
func (s *APIServer) Run() error {
	s.logger.Info().Str("port", s.config.Port).Msg("Starting ONVIF API server")
	return s.router.Run(":" + s.config.Port)
}
