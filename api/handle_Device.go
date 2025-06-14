package api

import (
	"go-onvif/api/handlers"
	"go-onvif/api/service"
	"go-onvif/api/utils"
	onvif "go-onvif/internal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	scanTimeout = 10 * time.Second
)

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
	results := s.deviceScanner.PerformDeviceScan(ips, req, acceptedData, scanTimeoutDuration)

	var aliveResults []utils.DeviceScanResult
	for _, r := range results {
		if r.Alive {
			aliveResults = append(aliveResults, r)
		}
	}

	aliveCount := utils.CountAliveDevices(aliveResults)

	s.logger.Info().
		Int("total_scanned", len(results)).
		Int("alive_devices", aliveCount).
		Msg("ONVIF device scan completed")

	// Return results
	c.JSON(http.StatusOK, gin.H{
		"results": aliveResults,
		"summary": gin.H{
			"total_scanned": len(results),
			"alive_devices": aliveCount,
		},
	})
}

func (s *APIServer) getDeviceCache(c *gin.Context) {
	type Services struct {
		Analytics string `json:"analytics,omitempty"`
		Device    string `json:"device,omitempty"`
		DeviceIO  string `json:"deviceio,omitempty"`
		Events    string `json:"events,omitempty"`
		Imaging   string `json:"imaging,omitempty"`
		Media     string `json:"media,omitempty"`
	}

	type CachedDeviceInfo struct {
		Xaddr    string   `json:"xaddr"`
		Username string   `json:"username"`
		Services Services `json:"services"`
	}

	cachedDevices := s.deviceCache.GetAllDevices()
	response := make(map[string]CachedDeviceInfo)

	for ip, dev := range cachedDevices {
		deviceParams := dev.GetDeviceParams()
		svcs := dev.GetServices()

		response[ip] = CachedDeviceInfo{
			Xaddr:    deviceParams.Xaddr,
			Username: deviceParams.Username,
			Services: Services{
				Analytics: svcs["analytics"],
				Device:    svcs["device"],
				DeviceIO:  svcs["deviceio"],
				Events:    svcs["events"],
				Imaging:   svcs["imaging"],
				Media:     svcs["media"],
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"cached_devices": response,
	})
}

func (s *APIServer) getDeviceDetails(c *gin.Context) {
	xaddr := c.Query("ip")
	username := c.Query("username")
	password := c.Query(("password"))
	dev, err := s.deviceCache.GetDevice(xaddr, username, password)
	deviceDetails := make(map[string]interface{})
	if err == nil {
		details := s.extractDeviceEssentials(xaddr, dev)
		deviceDetails[xaddr] = details
	}

	c.JSON(http.StatusOK, gin.H{
		"device_details": deviceDetails,
	})
}

func (s *APIServer) extractDeviceEssentials(ip string, dev *onvif.Device) map[string]interface{} {
	details := make(map[string]interface{})

	// Extract Device Information
	if devInfo := s.getDeviceInformation(ip, dev); devInfo != nil {
		details["device_info"] = devInfo
	}

	// Extract Network Interfaces
	if netInfo := s.getNetworkInfo(ip, dev); netInfo != nil {
		details["network_info"] = netInfo
	}

	// Extract Media Profiles and Stream URIs
	if mediaInfo := handlers.GetMediaInfo(ip, dev); mediaInfo != nil {
		details["media_info"] = mediaInfo
	}

	return details
}

func (s *APIServer) getDeviceInformation(ip string, dev *onvif.Device) interface{} {
	devInfoResp, err := service.CallDeviceMethod("GetDeviceInformation", dev, nil)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to get device information")
		return map[string]interface{}{"error": err.Error()}
	}

	return devInfoResp
}

func (s *APIServer) getNetworkInfo(ip string, dev *onvif.Device) interface{} {
	netResp, err := service.CallDeviceMethod("GetNetworkInterfaces", dev, nil)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to get network interfaces")
		return map[string]interface{}{"error": err.Error()}
	}

	return netResp
}
