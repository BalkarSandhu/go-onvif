package extractors

import (
	onvif "go-onvif/internal"
)

func ExtractProfileInfo(dev *onvif.Device, profileMap map[string]interface{}) ProfileInfo {
	info := ProfileInfo{}

	if name, ok := profileMap["Name"].(string); ok {
		info.Name = name
	}
	if token, ok := profileMap["Token"].(string); ok {

		stream_uri, err := GetStreamURI(dev, token)

		if err == nil {
			info.StreamUri = stream_uri
		}

		info.Token = token
	}
	if fixed, ok := profileMap["Fixed"].(bool); ok {
		info.Fixed = fixed
	}

	if videoEncoder, exists := profileMap["VideoEncoderConfiguration"]; exists {
		info.VideoEncoder = ExtractVideoEncoderInfo(videoEncoder)
	}
	if videoSource, exists := profileMap["VideoSourceConfiguration"]; exists {
		info.VideoSource = ExtractVideoSourceInfo(videoSource)
	}
	if analytics, exists := profileMap["VideoAnalyticsConfiguration"]; exists {
		info.Analytics = ExtractAnalyticsInfo(analytics)
	}
	if ptz, exists := profileMap["PTZConfiguration"]; exists {
		info.PTZ = ExtractPTZInfo(ptz)
	}
	if metadata, exists := profileMap["MetadataConfiguration"]; exists {
		info.Metadata = ExtractMetadataInfo(metadata)
	}

	return info
}

// extractVideoEncoderInfo extracts video encoder configuration
func ExtractVideoEncoderInfo(videoEncoder interface{}) *VideoEncoderInfo {
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
		info.RateControl = ExtractRateControlInfo(rateControl)
	}

	// Extract H264 configuration
	if h264, exists := videoEncoderMap["H264"]; exists {
		info.H264 = ExtractH264Info(h264)
	}

	// Extract multicast info
	if multicast, exists := videoEncoderMap["Multicast"]; exists {
		info.Multicast = ExtractMulticastInfo(multicast)
	}

	return info
}

// extractRateControlInfo extracts rate control information
func ExtractRateControlInfo(rateControl interface{}) *RateControlInfo {
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
func ExtractH264Info(h264 interface{}) *H264Info {
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
func ExtractMulticastInfo(multicast interface{}) *MulticastInfo {
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
func ExtractVideoSourceInfo(videoSource interface{}) *VideoSourceInfo {
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
func ExtractAnalyticsInfo(analytics interface{}) *AnalyticsInfo {
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
func ExtractPTZInfo(ptz any) any {
	return ptz
}

// extractMetadataInfo extracts metadata configuration
func ExtractMetadataInfo(metadata interface{}) *MetadataInfo {
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
