package service

import (
	"go-onvif/internal/xsd/onvif"
	"math"
)

// PTZRequest represents the API request payload
type AbsoluteMove2 struct {
	ProfileToken onvif.ReferenceToken
	Position     PTZVector
}

type GeoMove2 struct {
	ProfileToken    onvif.ReferenceToken
	TargetLatitude  float64
	TargetLongitude float64
	SelfLatitude    float64
	SelfLongitude   float64
	SelfHeading     uint32
	CameraHeight    float64
	Zoom            float64
}

type PTZVector struct {
	PanTilt Vector2D
	Zoom    Vector1D
}

type Vector2D struct {
	X float64
	Y float64
}

type Vector1D struct {
	X float64
}

type Target struct {
	X float64
	Y float64
}

// Helper function to convert degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// Helper function to get meters per degree of latitude
func metersPerDegreeLatitude() float64 {
	// Approximately 111,320 meters per degree of latitude (constant everywhere)
	return 111320.0
}

// Helper function to get meters per degree of longitude at a given latitude
func metersPerDegreeLongitude(latitude float64) float64 {
	// Meters per degree of longitude varies with latitude
	// Formula: 111,320 × cos(latitude)
	return 111320.0 * math.Cos(degreesToRadians(latitude))
}

// Main conversion function
// cameraHeight: height of camera/sensor above ground in meters (e.g., 10 meters for tower height)
func convertLatLongToXY(targetLatitude, targetLongitude float64, selfHeading uint32, selfLatitude, selfLongitude float64, cameraHeight float64) Target {
	// Step 1: Calculate the difference in latitude and longitude
	latDiff := targetLatitude - selfLatitude
	longDiff := targetLongitude - selfLongitude

	// Step 2: Convert differences to meters
	yMeters := latDiff * metersPerDegreeLatitude()
	xMeters := longDiff * metersPerDegreeLongitude(selfLatitude)

	// Step 3: Convert the heading from degrees to radians
	theta := degreesToRadians(float64(selfHeading))

	// Step 4: Rotate coordinates from world frame to vehicle/camera frame
	// Apply inverse rotation to align with vehicle's heading
	rotatedX := xMeters*math.Cos(-theta) + yMeters*math.Sin(-theta)
	rotatedY := yMeters*math.Cos(-theta) - xMeters*math.Sin(-theta)

	// Step 5: Calculate the pan angle (horizontal angle from forward direction)
	// atan2(x, y) gives angle in radians from forward axis
	// Result range: [-π, π] where 0 is forward, π/2 is right, -π/2 is left
	panAngleRad := math.Atan2(rotatedX, rotatedY)

	// Step 6: Calculate horizontal distance and tilt angle
	horizontalDistance := math.Sqrt(rotatedX*rotatedX + rotatedY*rotatedY)

	// Calculate height difference
	// Since target is on ground and camera is elevated, height difference is negative
	// (looking down from camera to target)
	heightDiff := -cameraHeight

	// atan2(height, distance) gives elevation angle
	// Result range: [-π/2, π/2] where 0 is level, positive is up, negative is down
	// With negative heightDiff, this will give a negative angle (looking down)
	tiltAngleRad := math.Atan2(heightDiff, horizontalDistance)

	// Step 7: Convert angles to your coordinate system [-2, 0]
	// Pan mapping: 0° → -0.5, 90° → -1.0, 180° → -1.5, -90° → 0.0
	targetX := -0.5 - (panAngleRad / math.Pi)

	// Tilt mapping: 0° → -0.5, -90° → -1.0, 90° → 0.0, 180° → -1.5
	// Note: We add (not subtract) because positive angle (up) should move toward 0
	targetY := -0.5 + (tiltAngleRad / (math.Pi / 2))

	// Step 8: Wrap X to [-2, 0] range
	if targetX < -2.0 {
		targetX += 2.0
	} else if targetX > 0.0 {
		targetX -= 2.0
	}

	// Step 9: Wrap Y to [-2, 0] range
	if targetY < -2.0 {
		targetY += 2.0
	} else if targetY > 0.0 {
		targetY -= 2.0
	}

	return Target{
		X: targetX,
		Y: targetY,
	}
}
