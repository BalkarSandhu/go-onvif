package api

import (
	"os"
	"time"

	"go-onvif/json_apis/cache"
	"go-onvif/json_apis/config"
	"go-onvif/json_apis/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type APIServer struct {
	router       *gin.Engine
	logger       zerolog.Logger
	deviceCache  *cache.DeviceCache
	devceScanner utils.DeviceScanner
	limiter      *rate.Limiter
	config       config.Config
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
