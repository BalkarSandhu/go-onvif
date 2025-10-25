package service

import (
	"encoding/json"
	"fmt"
	"go-onvif/internal/ptz"
	"go-onvif/internal/xsd"
	"go-onvif/internal/xsd/onvif"
	"math"
)

// PTZGeoController handles PTZ camera geo-targeting with deterministic behavior
type PTZGeoController struct {
	ProfileToken      string
	MaxVelocity       float64
	VelocityScale     float64
	Timeout           string
	MinZoomDistance   float64
	MaxZoomDistance   float64
	CurrentZoom       float64
	ZoomVelocityScale float64
	DeadbandPan       float64 // Degrees - stops movement within this threshold
	DeadbandTilt      float64 // Degrees
	DeadbandZoom      float64 // Normalized units
	Precision         int     // Decimal places for rounding (default: 6)
}

// PTZRequest represents the API request payload
type PTZRequest struct {
	ProfileToken string   `json:"ProfileToken"`
	Camera       Camera   `json:"camera"`
	Target       Target   `json:"target"`
	Settings     Settings `json:"settings,omitempty"`
}

// Camera represents camera position and state
type Camera struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Height     float64 `json:"height"`
	YawOffset  float64 `json:"yaw_offset"`  // Camera mounting yaw offset in degrees (positive = rotated clockwise)
	TiltOffset float64 `json:"tilt_offset"` // Camera mounting tilt offset in degrees (positive = tilted up from horizon)
	PanRange   struct {
		Min float64 `json:"min"` // Minimum pan angle in degrees
		Max float64 `json:"max"` // Maximum pan angle in degrees
	} `json:"pan_range"`
	TiltRange struct {
		Min float64 `json:"min"` // Minimum tilt angle in degrees (negative = down)
		Max float64 `json:"max"` // Maximum tilt angle in degrees (positive = up)
	} `json:"tilt_range"`
}

// Target represents target coordinates
type Target struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"`
}

// Settings represents optional PTZ settings
type Settings struct {
	MaxVelocity       *float64 `json:"max_velocity,omitempty"`
	VelocityScale     *float64 `json:"velocity_scale,omitempty"`
	Timeout           *string  `json:"timeout,omitempty"`
	MinZoomDistance   *float64 `json:"min_zoom_distance,omitempty"`
	MaxZoomDistance   *float64 `json:"max_zoom_distance,omitempty"`
	ZoomVelocityScale *float64 `json:"zoom_velocity_scale,omitempty"`
	DeadbandPan       *float64 `json:"deadband_pan,omitempty"`
	DeadbandTilt      *float64 `json:"deadband_tilt,omitempty"`
	DeadbandZoom      *float64 `json:"deadband_zoom,omitempty"`
	Precision         *int     `json:"precision,omitempty"`
}

// NewPTZGeoController creates a new controller with default settings
func NewPTZGeoController(profileToken string) *PTZGeoController {
	return &PTZGeoController{
		ProfileToken:      profileToken,
		MaxVelocity:       1.0,
		VelocityScale:     0.5,
		Timeout:           "PT0.5S",
		MinZoomDistance:   10.0,
		MaxZoomDistance:   500.0,
		CurrentZoom:       0.0,
		ZoomVelocityScale: 0.3,
		DeadbandPan:       0.5,  // 0.5 degrees
		DeadbandTilt:      0.3,  // 0.3 degrees
		DeadbandZoom:      0.02, // 2% of zoom range
		Precision:         6,    // 6 decimal places for deterministic output
	}
}

// ApplySettings applies custom settings to the controller
func (c *PTZGeoController) ApplySettings(settings Settings) {
	if settings.MaxVelocity != nil {
		c.MaxVelocity = *settings.MaxVelocity
	}
	if settings.VelocityScale != nil {
		c.VelocityScale = *settings.VelocityScale
	}
	if settings.Timeout != nil {
		c.Timeout = *settings.Timeout
	}
	if settings.MinZoomDistance != nil {
		c.MinZoomDistance = *settings.MinZoomDistance
	}
	if settings.MaxZoomDistance != nil {
		c.MaxZoomDistance = *settings.MaxZoomDistance
	}
	if settings.ZoomVelocityScale != nil {
		c.ZoomVelocityScale = *settings.ZoomVelocityScale
	}
	if settings.DeadbandPan != nil {
		c.DeadbandPan = *settings.DeadbandPan
	}
	if settings.DeadbandTilt != nil {
		c.DeadbandTilt = *settings.DeadbandTilt
	}
	if settings.DeadbandZoom != nil {
		c.DeadbandZoom = *settings.DeadbandZoom
	}
	if settings.Precision != nil {
		c.Precision = *settings.Precision
	}
}

// roundToPrecision rounds value to specified decimal places for deterministic output
func (c *PTZGeoController) roundToPrecision(val float64) float64 {
	multiplier := math.Pow(10, float64(c.Precision))
	return math.Round(val*multiplier) / multiplier
}

// CalculateBearingAndDistance calculates bearing and distance between two coordinates
// Returns bearing in degrees (0-360) and distance in meters
// Results are deterministic for same inputs
func (c *PTZGeoController) CalculateBearingAndDistance(camLat, camLon, targetLat, targetLon float64) (bearing, distance float64, err error) {
	// Validate coordinates
	if camLat < -90 || camLat > 90 || targetLat < -90 || targetLat > 90 {
		return 0, 0, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	if camLon < -180 || camLon > 180 || targetLon < -180 || targetLon > 180 {
		return 0, 0, fmt.Errorf("invalid longitude: must be between -180 and 180")
	}

	// Convert to radians
	lat1 := camLat * math.Pi / 180.0
	lon1 := camLon * math.Pi / 180.0
	lat2 := targetLat * math.Pi / 180.0
	lon2 := targetLon * math.Pi / 180.0

	dlon := lon2 - lon1

	// Calculate bearing
	x := math.Sin(dlon) * math.Cos(lat2)
	y := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)

	bearing = math.Atan2(x, y) * 180.0 / math.Pi
	bearing = math.Mod(bearing+360.0, 360.0)

	// Calculate distance using Haversine formula
	dlat := lat2 - lat1
	a := math.Sin(dlat/2.0)*math.Sin(dlat/2.0) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2.0)*math.Sin(dlon/2.0)
	circum := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	distance = 6371000.0 * circum // Earth radius in meters

	// Round for deterministic output
	bearing = c.roundToPrecision(bearing)
	distance = c.roundToPrecision(distance)

	return bearing, distance, nil
}

// CalculateElevationAngle calculates tilt angle in degrees
// Positive angle = tilt up, negative = tilt down
// heightDiff = camHeight - targetHeight
func (c *PTZGeoController) CalculateElevationAngle(distance, heightDiff float64) float64 {
	if distance == 0.0 {
		if heightDiff > 0.0 {
			return -90.0 // Look straight down
		} else if heightDiff < 0.0 {
			return 90.0 // Look straight up
		}
		return 0.0
	}

	// heightDiff = camHeight - targetHeight
	// Positive heightDiff (camera above target) -> negative angle (tilt down)
	// Negative heightDiff (camera below target) -> positive angle (tilt up)
	elevation := math.Atan2(heightDiff, distance) * 180.0 / math.Pi

	// Invert: we want negative for down, positive for up
	elevation = -elevation

	return c.roundToPrecision(elevation)
}

// CalculateTargetZoom calculates the target zoom level based on distance
// Returns normalized zoom value (0 = wide, 1 = full zoom)
func (c *PTZGeoController) CalculateTargetZoom(distance float64) float64 {
	if distance <= c.MinZoomDistance {
		return 1.0
	}
	if distance >= c.MaxZoomDistance {
		return 0.0
	}
	targetZoom := 1.0 - (distance-c.MinZoomDistance)/(c.MaxZoomDistance-c.MinZoomDistance)
	return c.roundToPrecision(targetZoom)
}

// NormalizedToAngle converts ONVIF normalized position (-1..1) to angle in degrees
func NormalizedToAngle(normalized, minAngle, maxAngle float64) float64 {
	// ONVIF: -1 = min, 0 = center, +1 = max
	return minAngle + (normalized+1.0)/2.0*(maxAngle-minAngle)
}

// AngleToNormalized converts angle in degrees to ONVIF normalized position (-1..1)
func AngleToNormalized(angle, minAngle, maxAngle float64) float64 {
	if maxAngle == minAngle {
		return 0.0
	}
	normalized := 2.0*(angle-minAngle)/(maxAngle-minAngle) - 1.0
	return clamp(normalized, -1.0, 1.0)
}

// NormalizeAngle wraps angle to [-180, 180) range
func NormalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360.0)
	if angle >= 180.0 {
		angle -= 360.0
	} else if angle < -180.0 {
		angle += 360.0
	}
	return angle
}

// Clamp value between min and max
func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// applyDeadband returns 0 if error is within deadband, otherwise returns error
func applyDeadband(error, deadband float64) float64 {
	if math.Abs(error) < deadband {
		return 0.0
	}
	return error
}

// GetMovementCommand generates PTZ movement command to focus on target coordinates
func (c *PTZGeoController) GetMovementCommand(
	camLat, camLon, camHeight float64,
	currentPanNorm, currentTiltNorm float64, // Normalized ONVIF positions (-1..1)
	targetLat, targetLon, targetHeight float64,
	panRangeMin, panRangeMax, tiltRangeMin, tiltRangeMax float64,
	yawOffset, tiltOffset float64) (*ptz.ContinuousMove, error) {

	// Handle cameras without absolute position support
	if panRangeMin == 0.0 && panRangeMax == 0.0 {
		return c.GetMovementCommandNormalized(
			camLat, camLon, camHeight,
			currentPanNorm, currentTiltNorm,
			targetLat, targetLon, targetHeight,
			yawOffset, tiltOffset)
	}

	// Step 1: Convert normalized positions to angles in degrees
	currentPanDeg := NormalizedToAngle(currentPanNorm, panRangeMin, panRangeMax)
	currentTiltDeg := NormalizedToAngle(currentTiltNorm, tiltRangeMin, tiltRangeMax)

	currentPanDeg = c.roundToPrecision(currentPanDeg)
	currentTiltDeg = c.roundToPrecision(currentTiltDeg)

	// Step 2: Compute target bearing and elevation in world coordinates
	targetBearing, distance, err := c.CalculateBearingAndDistance(camLat, camLon, targetLat, targetLon)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate bearing: %w", err)
	}

	heightDiff := camHeight - targetHeight
	targetElevation := c.CalculateElevationAngle(distance, heightDiff)

	// Step 3: Apply mounting offsets to convert world coordinates to camera coordinates
	// yawOffset: positive = camera rotated clockwise from north
	// tiltOffset: positive = camera tilted up from horizon
	targetBearingCam := targetBearing + yawOffset
	targetElevationCam := targetElevation + tiltOffset

	// Step 4: Clamp target to physical limits BEFORE computing differences
	targetBearingCam = clamp(targetBearingCam, panRangeMin, panRangeMax)
	targetElevationCam = clamp(targetElevationCam, tiltRangeMin, tiltRangeMax)

	targetBearingCam = c.roundToPrecision(targetBearingCam)
	targetElevationCam = c.roundToPrecision(targetElevationCam)

	// Step 5: Compute shortest angular difference (handles wraparound for pan)
	panDiff := NormalizeAngle(targetBearingCam - currentPanDeg)
	tiltDiff := targetElevationCam - currentTiltDeg

	panDiff = c.roundToPrecision(panDiff)
	tiltDiff = c.roundToPrecision(tiltDiff)

	// Step 6: Apply deadband to prevent oscillation
	panDiff = applyDeadband(panDiff, c.DeadbandPan)
	tiltDiff = applyDeadband(tiltDiff, c.DeadbandTilt)

	// Step 7: Compute velocities (proportional to angular error)
	// Scale by velocity parameters
	panVelocity := clamp((panDiff/180.0)*c.VelocityScale, -c.MaxVelocity, c.MaxVelocity)
	tiltVelocity := clamp((tiltDiff/90.0)*c.VelocityScale, -c.MaxVelocity, c.MaxVelocity)

	// Step 8: Compute zoom velocity
	targetZoom := c.CalculateTargetZoom(distance)
	zoomDiff := targetZoom - c.CurrentZoom
	zoomDiff = c.roundToPrecision(zoomDiff)
	zoomDiff = applyDeadband(zoomDiff, c.DeadbandZoom)

	zoomVelocity := clamp(zoomDiff*c.ZoomVelocityScale, -c.MaxVelocity, c.MaxVelocity)

	// Step 9: Round velocities for deterministic output
	panVelocity = c.roundToPrecision(panVelocity)
	tiltVelocity = c.roundToPrecision(tiltVelocity)
	zoomVelocity = c.roundToPrecision(zoomVelocity)

	// Step 10: Build ONVIF PTZ command
	command := &ptz.ContinuousMove{
		ProfileToken: onvif.ReferenceToken(c.ProfileToken),
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{
				X: panVelocity,
				Y: tiltVelocity,
			},
			Zoom: onvif.Vector1D{
				X: zoomVelocity,
			},
		},
		Timeout: xsd.Duration(c.Timeout),
	}

	return command, nil
}

// GetMovementCommandNormalized handles cameras without absolute position support
// Uses simplified velocity-based approach without angle conversion
func (c *PTZGeoController) GetMovementCommandNormalized(
	camLat, camLon, camHeight float64,
	currentPanNorm, currentTiltNorm float64,
	targetLat, targetLon, targetHeight float64,
	yawOffset, tiltOffset float64) (*ptz.ContinuousMove, error) {

	// Calculate target bearing and distance
	targetBearing, distance, err := c.CalculateBearingAndDistance(camLat, camLon, targetLat, targetLon)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate bearing: %w", err)
	}

	heightDiff := camHeight - targetHeight
	targetElevation := c.CalculateElevationAngle(distance, heightDiff)

	// For cameras without absolute positioning, use assumed ranges
	// These should be calibrated for each camera model
	const (
		assumedPanRange  = 340.0 // Most PTZ cameras have ~340° rotation
		assumedTiltRange = 180.0 // -90 to +90 degrees
	)

	// Convert normalized to approximate degrees
	currentPanDeg := currentPanNorm * (assumedPanRange / 2.0)
	currentTiltDeg := currentTiltNorm * (assumedTiltRange / 2.0)

	currentPanDeg = c.roundToPrecision(currentPanDeg)
	currentTiltDeg = c.roundToPrecision(currentTiltDeg)

	// Apply mounting offsets
	targetBearingCam := targetBearing + yawOffset
	targetElevationCam := targetElevation + tiltOffset

	// Normalize target bearing to assumed range
	targetBearingCam = NormalizeAngle(targetBearingCam)

	// Clamp to assumed pan range
	targetBearingCam = clamp(targetBearingCam, -(assumedPanRange / 2.0), (assumedPanRange / 2.0))

	// Clamp tilt to assumed range
	targetElevationCam = clamp(targetElevationCam, -90.0, 90.0)

	targetBearingCam = c.roundToPrecision(targetBearingCam)
	targetElevationCam = c.roundToPrecision(targetElevationCam)

	// Calculate differences
	panDiff := NormalizeAngle(targetBearingCam - currentPanDeg)
	tiltDiff := targetElevationCam - currentTiltDeg

	panDiff = c.roundToPrecision(panDiff)
	tiltDiff = c.roundToPrecision(tiltDiff)

	// Apply deadband
	panDiff = applyDeadband(panDiff, c.DeadbandPan)
	tiltDiff = applyDeadband(tiltDiff, c.DeadbandTilt)

	// Compute velocities with conservative scaling for uncalibrated cameras
	conservativeScale := c.VelocityScale * 0.7 // Reduce speed by 30% for safety
	panVelocity := clamp((panDiff/180.0)*conservativeScale, -c.MaxVelocity, c.MaxVelocity)
	tiltVelocity := clamp((tiltDiff/90.0)*conservativeScale, -c.MaxVelocity, c.MaxVelocity)

	// Compute zoom
	targetZoom := c.CalculateTargetZoom(distance)
	zoomDiff := targetZoom - c.CurrentZoom
	zoomDiff = c.roundToPrecision(zoomDiff)
	zoomDiff = applyDeadband(zoomDiff, c.DeadbandZoom)

	zoomVelocity := clamp(zoomDiff*c.ZoomVelocityScale, -c.MaxVelocity, c.MaxVelocity)

	// Round velocities for deterministic output
	panVelocity = c.roundToPrecision(panVelocity)
	tiltVelocity = c.roundToPrecision(tiltVelocity)
	zoomVelocity = c.roundToPrecision(zoomVelocity)

	// Build command
	command := &ptz.ContinuousMove{
		ProfileToken: onvif.ReferenceToken(c.ProfileToken),
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{
				X: panVelocity,
				Y: tiltVelocity,
			},
			Zoom: onvif.Vector1D{
				X: zoomVelocity,
			},
		},
		Timeout: xsd.Duration(c.Timeout),
	}

	return command, nil
}

// GetMoveParams generates PTZ movement parameters from request and current status
func GetMoveParams(status ptz.GetStatusResponse, data []byte) (*ptz.ContinuousMove, error) {
	// 1️⃣ Unmarshal request
	var request PTZRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if request.ProfileToken == "" {
		return nil, fmt.Errorf("missing ProfileToken")
	}

	// 2️⃣ Validate coordinates
	if request.Camera.Latitude < -90 || request.Camera.Latitude > 90 {
		return nil, fmt.Errorf("invalid camera latitude: %f", request.Camera.Latitude)
	}
	if request.Camera.Longitude < -180 || request.Camera.Longitude > 180 {
		return nil, fmt.Errorf("invalid camera longitude: %f", request.Camera.Longitude)
	}
	if request.Target.Latitude < -90 || request.Target.Latitude > 90 {
		return nil, fmt.Errorf("invalid target latitude: %f", request.Target.Latitude)
	}
	if request.Target.Longitude < -180 || request.Target.Longitude > 180 {
		return nil, fmt.Errorf("invalid target longitude: %f", request.Target.Longitude)
	}

	// 3️⃣ Check for absolute positioning support
	hasAbsolutePositioning := !(request.Camera.PanRange.Min == 0.0 && request.Camera.PanRange.Max == 0.0)

	if !hasAbsolutePositioning {
		fmt.Println("Warning: Camera doesn't support absolute positioning. Using approximation mode. Calibration recommended.")
	}

	// 4️⃣ Extract current PTZ position from GetStatusResponse (normalized -1..1)
	currentPanNorm := float64(status.PTZStatus.Position.PanTilt.X)
	currentTiltNorm := float64(status.PTZStatus.Position.PanTilt.Y)
	currentZoom := float64(status.PTZStatus.Position.Zoom.X)

	// 5️⃣ Create controller
	controller := NewPTZGeoController(request.ProfileToken)
	controller.CurrentZoom = currentZoom
	controller.ApplySettings(request.Settings)

	// 6️⃣ Generate movement command
	cmd, err := controller.GetMovementCommand(
		request.Camera.Latitude, request.Camera.Longitude, request.Camera.Height,
		currentPanNorm, currentTiltNorm,
		request.Target.Latitude, request.Target.Longitude, request.Target.Height,
		request.Camera.PanRange.Min, request.Camera.PanRange.Max,
		request.Camera.TiltRange.Min, request.Camera.TiltRange.Max,
		request.Camera.YawOffset, request.Camera.TiltOffset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate movement command: %w", err)
	}

	return cmd, nil
}
