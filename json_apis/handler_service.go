package api

import (
	"net/http"
	"strings"

	"go-onvif/json_apis/utils"
	onvif "go-onvif/internal"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

func (s *APIServer) handleServiceMethod(c *gin.Context) {
	serviceName := c.Param("service")
	methodName := c.Param("method")
	username := c.GetHeader("username")
	password := c.GetHeader("password")
	xaddr := c.GetHeader("xaddr")

	if xaddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing xaddr header"})
		return
	}

	rawData, err := c.GetRawData()
	if err != nil {
		s.logger.Debug().Err(err).Msg("Failed to get raw data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	dev, err := s.deviceCache.GetDevice(xaddr, username, password)
	if err != nil {
		s.logger.Error().Err(err).Str("xaddr", xaddr).Msg("Failed to connect to device")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect to device: " + err.Error()})
		return
	}

	resp, err := s.CallServiceMethod(serviceName, methodName, rawData, dev)
	if err != nil {
		s.logger.Error().Err(err).Str("service", serviceName).Str("method", methodName).Msg("Method call failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *APIServer) CallServiceMethod(serviceName, methodName string, data []byte, dev *onvif.Device) (interface{}, error) {
	switch strings.ToLower(serviceName) {
	case "device":
		return utils.CallDeviceMethod(methodName, dev, data)
	case "ptz":
		return utils.CallPTZMethod(methodName, dev, data)
	case "media":
		return utils.CallMediaMethod(methodName, dev, data)
	default:
		return nil, errors.New("unknown service: " + serviceName)
	}
}
