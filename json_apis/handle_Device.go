package api

import (
	"encoding/json"
	"fmt"
	"go-onvif/json_apis/utils"
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
	results := s.devceScanner.PerformDeviceScan(ips, req, acceptedData, scanTimeoutDuration)

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
	cachedDevices := s.deviceCache.GetAllDevices()
	deviceDetails := make(map[string]interface{})

	for ip, dev := range cachedDevices {
		details := s.extractDeviceEssentials(ip, dev)
		deviceDetails[ip] = details
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
	if mediaInfo := s.getMediaInfo(ip, dev); mediaInfo != nil {
		details["media_info"] = mediaInfo
	}

	return details
}

func (s *APIServer) getDeviceInformation(ip string, dev *onvif.Device) interface{} {
	devInfoResp, err := utils.CallDeviceMethod("GetDeviceInformation", dev, nil)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to get device information")
		return map[string]interface{}{"error": err.Error()}
	}

	return devInfoResp
}

func (s *APIServer) getNetworkInfo(ip string, dev *onvif.Device) interface{} {
	netResp, err := utils.CallDeviceMethod("GetNetworkInterfaces", dev, nil)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to get network interfaces")
		return map[string]interface{}{"error": err.Error()}
	}

	return netResp
}

// ProfileInfo represents the extracted important information from a media profile
type ProfileInfo struct {
	Name         string            `json:"name,omitempty"`
	Token        string            `json:"token,omitempty"`
	Fixed        bool              `json:"fixed,omitempty"`
	VideoEncoder *VideoEncoderInfo `json:"video_encoder,omitempty"`
	VideoSource  *VideoSourceInfo  `json:"video_source,omitempty"`
	Analytics    *AnalyticsInfo    `json:"analytics,omitempty"`
	PTZ          *PTZInfo          `json:"ptz,omitempty"`
	Metadata     *MetadataInfo     `json:"metadata,omitempty"`
}

type VideoEncoderInfo struct {
	Encoding    string                 `json:"encoding,omitempty"`
	Quality     float64                `json:"quality,omitempty"`
	Resolution  map[string]interface{} `json:"resolution,omitempty"`
	RateControl *RateControlInfo       `json:"rate_control,omitempty"`
	H264        *H264Info              `json:"h264,omitempty"`
	Multicast   *MulticastInfo         `json:"multicast,omitempty"`
}

type RateControlInfo struct {
	BitrateLimit   int `json:"bitrate_limit,omitempty"`
	FrameRateLimit int `json:"frame_rate_limit,omitempty"`
}

type H264Info struct {
	Profile   string `json:"profile,omitempty"`
	GovLength int    `json:"gov_length,omitempty"`
}

type MulticastInfo struct {
	Address   string `json:"address,omitempty"`
	Port      int    `json:"port,omitempty"`
	AutoStart bool   `json:"auto_start,omitempty"`
}

type VideoSourceInfo struct {
	Bounds      map[string]interface{} `json:"bounds,omitempty"`
	SourceToken string                 `json:"source_token,omitempty"`
	ViewMode    string                 `json:"view_mode,omitempty"`
}

type AnalyticsInfo struct {
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

type PTZInfo struct {
	NodeToken    string                 `json:"node_token,omitempty"`
	DefaultSpeed map[string]interface{} `json:"default_speed,omitempty"`
}

type MetadataInfo struct {
	Analytics bool `json:"analytics,omitempty"`
	Events    bool `json:"events,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *APIServer) getMediaInfo(ip string, dev *onvif.Device) interface{} {
	// Get profiles from device
	profilesResp, err := utils.CallMediaMethod("GetProfiles", dev, nil)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to get media profiles")
		return ErrorResponse{Error: err.Error()}
	}

	// Parse the response
	profiles, err := s.parseProfilesResponse(profilesResp, ip)
	if err != nil {
		return ErrorResponse{Error: err.Error()}
	}

	// Extract important information from each profile
	importantProfiles := make([]ProfileInfo, 0, len(profiles))
	for i, profile := range profiles {
		profileInfo := s.extractProfileInfo(profile)
		importantProfiles = append(importantProfiles, profileInfo)
		s.logger.Debug().Int("profile_index", i).Str("ip", ip).Msg("Profile info extracted")
	}

	s.logger.Info().Int("profile_count", len(importantProfiles)).Str("ip", ip).Msg("Successfully extracted media profiles")
	return importantProfiles
}

// parseProfilesResponse parses the raw profiles response into a slice of maps
func (s *APIServer) parseProfilesResponse(profilesResp interface{}, ip string) ([]map[string]interface{}, error) {
	// Marshal and unmarshal to get consistent JSON structure
	jsonData, err := json.Marshal(profilesResp)
	if err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to marshal profiles response")
		return nil, fmt.Errorf("failed to marshal profiles response: %w", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		s.logger.Warn().Err(err).Str("ip", ip).Msg("Failed to unmarshal profiles JSON")
		return nil, fmt.Errorf("failed to unmarshal profiles JSON: %w", err)
	}

	// Extract profiles array
	profilesRaw, exists := parsed["Profiles"]
	if !exists {
		s.logger.Warn().Str("ip", ip).Msg("Profiles key not found in response")
		return nil, fmt.Errorf("profiles not found in response")
	}

	profilesSlice, ok := profilesRaw.([]interface{})
	if !ok {
		s.logger.Warn().Str("ip", ip).Msg("Profiles is not a slice")
		return nil, fmt.Errorf("profiles format invalid - expected array")
	}

	// Convert to slice of maps
	profiles := make([]map[string]interface{}, 0, len(profilesSlice))
	for i, profile := range profilesSlice {
		profileMap, ok := profile.(map[string]interface{})
		if !ok {
			s.logger.Warn().Int("profile_index", i).Str("ip", ip).Msg("Profile is not a map")
			continue // Skip invalid profiles instead of failing completely
		}
		profiles = append(profiles, profileMap)
	}

	return profiles, nil
}

// extractProfileInfo extracts important information from a single profile
func (s *APIServer) extractProfileInfo(profileMap map[string]interface{}) ProfileInfo {
	info := ProfileInfo{}

	// Basic profile info
	if name, ok := profileMap["Name"].(string); ok {
		info.Name = name
	}
	if token, ok := profileMap["Token"].(string); ok {
		info.Token = token
	}
	if fixed, ok := profileMap["Fixed"].(bool); ok {
		info.Fixed = fixed
	}

	// Extract video encoder configuration
	if videoEncoder, exists := profileMap["VideoEncoderConfiguration"]; exists {
		info.VideoEncoder = s.extractVideoEncoderInfo(videoEncoder)
	}

	// Extract video source configuration
	if videoSource, exists := profileMap["VideoSourceConfiguration"]; exists {
		info.VideoSource = s.extractVideoSourceInfo(videoSource)
	}

	// Extract analytics configuration
	if analytics, exists := profileMap["VideoAnalyticsConfiguration"]; exists {
		info.Analytics = s.extractAnalyticsInfo(analytics)
	}

	// Extract PTZ configuration
	if ptz, exists := profileMap["PTZConfiguration"]; exists {
		info.PTZ = s.extractPTZInfo(ptz)
	}

	// Extract metadata configuration
	if metadata, exists := profileMap["MetadataConfiguration"]; exists {
		info.Metadata = s.extractMetadataInfo(metadata)
	}

	return info
}

// extractVideoEncoderInfo extracts video encoder configuration
func (s *APIServer) extractVideoEncoderInfo(videoEncoder interface{}) *VideoEncoderInfo {
	videoEncoderMap, ok := videoEncoder.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &VideoEncoderInfo{}

	if encoding, ok := videoEncoderMap["Encoding"].(string); ok {
		info.Encoding = encoding
	}
	if quality, ok := videoEncoderMap["Quality"].(float64); ok {
		info.Quality = quality
	}
	if resolution, ok := videoEncoderMap["Resolution"].(map[string]interface{}); ok {
		info.Resolution = resolution
	}

	// Extract rate control
	if rateControl, exists := videoEncoderMap["RateControl"]; exists {
		info.RateControl = s.extractRateControlInfo(rateControl)
	}

	// Extract H264 configuration
	if h264, exists := videoEncoderMap["H264"]; exists {
		info.H264 = s.extractH264Info(h264)
	}

	// Extract multicast info
	if multicast, exists := videoEncoderMap["Multicast"]; exists {
		info.Multicast = s.extractMulticastInfo(multicast)
	}

	return info
}

// extractRateControlInfo extracts rate control information
func (s *APIServer) extractRateControlInfo(rateControl interface{}) *RateControlInfo {
	rateControlMap, ok := rateControl.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &RateControlInfo{}
	if bitrate, ok := rateControlMap["BitrateLimit"].(float64); ok {
		info.BitrateLimit = int(bitrate)
	}
	if frameRate, ok := rateControlMap["FrameRateLimit"].(float64); ok {
		info.FrameRateLimit = int(frameRate)
	}

	return info
}

// extractH264Info extracts H264 configuration
func (s *APIServer) extractH264Info(h264 interface{}) *H264Info {
	h264Map, ok := h264.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &H264Info{}
	if profile, ok := h264Map["H264Profile"].(string); ok {
		info.Profile = profile
	}
	if govLength, ok := h264Map["GovLength"].(float64); ok {
		info.GovLength = int(govLength)
	}

	return info
}

// extractMulticastInfo extracts multicast information
func (s *APIServer) extractMulticastInfo(multicast interface{}) *MulticastInfo {
	multicastMap, ok := multicast.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &MulticastInfo{}
	if address, ok := multicastMap["Address"].(string); ok {
		info.Address = address
	}
	if port, ok := multicastMap["Port"].(float64); ok {
		info.Port = int(port)
	}
	if autoStart, ok := multicastMap["AutoStart"].(bool); ok {
		info.AutoStart = autoStart
	}

	return info
}

// extractVideoSourceInfo extracts video source configuration
func (s *APIServer) extractVideoSourceInfo(videoSource interface{}) *VideoSourceInfo {
	videoSourceMap, ok := videoSource.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &VideoSourceInfo{}
	if bounds, ok := videoSourceMap["Bounds"].(map[string]interface{}); ok {
		info.Bounds = bounds
	}
	if sourceToken, ok := videoSourceMap["SourceToken"].(string); ok {
		info.SourceToken = sourceToken
	}
	if viewMode, ok := videoSourceMap["ViewMode"].(string); ok {
		info.ViewMode = viewMode
	}

	return info
}

// extractAnalyticsInfo extracts analytics configuration
func (s *APIServer) extractAnalyticsInfo(analytics interface{}) *AnalyticsInfo {
	analyticsMap, ok := analytics.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &AnalyticsInfo{}
	if name, ok := analyticsMap["Name"].(string); ok {
		info.Name = name
	}
	if token, ok := analyticsMap["Token"].(string); ok {
		info.Token = token
	}

	return info
}

// extractPTZInfo extracts PTZ configuration
func (s *APIServer) extractPTZInfo(ptz interface{}) *PTZInfo {
	ptzMap, ok := ptz.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &PTZInfo{}
	if nodeToken, ok := ptzMap["NodeToken"].(string); ok {
		info.NodeToken = nodeToken
	}
	if defaultSpeed, ok := ptzMap["DefaultPTZSpeed"].(map[string]interface{}); ok {
		info.DefaultSpeed = defaultSpeed
	}

	return info
}

// extractMetadataInfo extracts metadata configuration
func (s *APIServer) extractMetadataInfo(metadata interface{}) *MetadataInfo {
	metadataMap, ok := metadata.(map[string]interface{})
	if !ok {
		return nil
	}

	info := &MetadataInfo{}
	if analytics, ok := metadataMap["Analytics"].(bool); ok {
		info.Analytics = analytics
	}
	if events, ok := metadataMap["Events"].(bool); ok {
		info.Events = events
	}

	return info
}

// Run starts the API server
func (s *APIServer) Run() error {
	s.logger.Info().Str("port", s.config.Port).Msg("Starting ONVIF API server")
	return s.router.Run(":" + s.config.Port)
}
