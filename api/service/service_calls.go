package service

import (
	"net/http"
	"strings"

	"go-onvif/api/cache"
	onvif "go-onvif/internal"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

type Service struct {
	logger      zerolog.Logger
	deviceCache *cache.DeviceCache
}

func NewService(dc *cache.DeviceCache) Service {
	return Service{
		deviceCache: dc,
	}
}

// handleServiceMethod handles dynamic ONVIF service calls.
func (s *Service) HandleServiceMethod(c *gin.Context) {
	serviceName := c.Param("service")
	methodName := c.Param("method")

	xaddr := c.GetHeader("xaddr")
	username := c.GetHeader("username")
	password := c.GetHeader("password")

	if xaddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing xaddr header"})
		return
	}

	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	dev, err := s.deviceCache.GetDevice(xaddr, username, password)
	if err != nil {
		s.logger.Error().Err(err).Str("xaddr", xaddr).Msg("Failed to connect to device")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect to device: " + err.Error()})
		return
	}

	resp, err := s.callServiceMethod(serviceName, methodName, rawData, dev)
	if err != nil {
		s.logger.Error().Err(err).
			Str("service", serviceName).
			Str("method", methodName).
			Msg("Method call failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CallServiceMethod dispatches ONVIF method calls based on the service name.
func (s *Service) callServiceMethod(serviceName, methodName string, data []byte, dev *onvif.Device) (interface{}, error) {
	switch strings.ToLower(serviceName) {
	case "device":
		return CallDeviceMethod(methodName, dev, data)
	case "ptz":
		return CallPTZMethod(methodName, dev, data)
	case "media":
		return CallMediaMethod(methodName, dev, data)
	default:
		return nil, errors.New("unknown service: " + serviceName)
	}
}
