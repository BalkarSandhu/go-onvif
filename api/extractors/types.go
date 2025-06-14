package extractors

type ProfileInfo struct {
	Name         string            `json:"name,omitempty"`
	Token        string            `json:"token,omitempty"`
	StreamUri    string      				`json:"stream_uri,omitempty"`
	Fixed        bool              `json:"fixed,omitempty"`
	VideoEncoder *VideoEncoderInfo `json:"video_encoder,omitempty"`
	VideoSource  *VideoSourceInfo  `json:"video_source,omitempty"`
	Analytics    *AnalyticsInfo    `json:"analytics,omitempty"`
	PTZ          interface{}        `json:"ptz,omitempty"`
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

type ErrorResponse struct {
	Error string `json:"error"`
}
