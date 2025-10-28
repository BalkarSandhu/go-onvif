package main

import (
	"fmt"
	"math"
)

// PTZRequest represents the API request payload
type AbsoluteMove2 struct {
	ProfileToken string
	Position     PTZVector
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

// GeoMove2 request structure
type GeoMove2Request struct {
	ProfileToken    string  `json:"profileToken"`
	TargetLatitude  float64 `json:"targetLatitude"`
	TargetLongitude float64 `json:"targetLongitude"`
	SelfLatitude    float64 `json:"selfLatitude"`
	SelfLongitude   float64 `json:"selfLongitude"`
	SelfHeading     uint32  `json:"selfHeading"`
	CameraHeight    float64 `json:"cameraHeight"`
	Zoom            float64 `json:"zoom"`
}

// Helper function to convert degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// Helper function to get meters per degree of latitude
func metersPerDegreeLatitude() float64 {
	return 111320.0
}

// Helper function to get meters per degree of longitude at a given latitude
func metersPerDegreeLongitude(latitude float64) float64 {
	return 111320.0 * math.Cos(degreesToRadians(latitude))
}

// Main conversion function
func convertLatLongToXY(targetLatitude, targetLongitude float64, selfHeading uint32, selfLatitude, selfLongitude float64, cameraHeight float64) (targetX, targetY float64) {
	fmt.Println("\n========== GeoMove2 Conversion Process ==========")

	// Step 1: Calculate the difference in latitude and longitude
	latDiff := targetLatitude - selfLatitude
	longDiff := targetLongitude - selfLongitude
	fmt.Printf("Step 1 - Coordinate Differences:\n")
	fmt.Printf("  Latitude Diff:  %.9f degrees\n", latDiff)
	fmt.Printf("  Longitude Diff: %.9f degrees\n", longDiff)

	// Step 2: Convert differences to meters
	yMeters := latDiff * metersPerDegreeLatitude()
	xMeters := longDiff * metersPerDegreeLongitude(selfLatitude)
	fmt.Printf("\nStep 2 - Convert to Meters:\n")
	fmt.Printf("  X (East-West):  %.2f meters\n", xMeters)
	fmt.Printf("  Y (North-South): %.2f meters\n", yMeters)
	fmt.Printf("  Straight-line Distance: %.2f meters\n", math.Sqrt(xMeters*xMeters+yMeters*yMeters))

	// Step 3: Convert the heading from degrees to radians
	theta := degreesToRadians(float64(selfHeading))
	fmt.Printf("\nStep 3 - Heading Conversion:\n")
	fmt.Printf("  Self Heading: %d degrees (%.4f radians)\n", selfHeading, theta)

	// Step 4: Rotate coordinates from world frame to vehicle/camera frame
	rotatedX := xMeters*math.Cos(-theta) + yMeters*math.Sin(-theta)
	rotatedY := yMeters*math.Cos(-theta) - xMeters*math.Sin(-theta)
	fmt.Printf("\nStep 4 - Rotate to Camera Frame:\n")
	fmt.Printf("  Rotated X (right/left):    %.2f meters\n", rotatedX)
	fmt.Printf("  Rotated Y (forward/back):  %.2f meters\n", rotatedY)

	// Step 5: Calculate the pan angle
	panAngleRad := math.Atan2(rotatedX, rotatedY)
	panAngleDeg := panAngleRad * 180 / math.Pi
	fmt.Printf("\nStep 5 - Calculate Pan Angle:\n")
	fmt.Printf("  Pan Angle: %.2f degrees (%.4f radians)\n", panAngleDeg, panAngleRad)
	if panAngleDeg > 0 {
		fmt.Printf("  Direction: %.2f degrees to the RIGHT of forward\n", panAngleDeg)
	} else {
		fmt.Printf("  Direction: %.2f degrees to the LEFT of forward\n", -panAngleDeg)
	}

	// Step 6: Calculate horizontal distance and tilt angle
	horizontalDistance := math.Sqrt(rotatedX*rotatedX + rotatedY*rotatedY)
	heightDiff := -cameraHeight
	tiltAngleRad := math.Atan2(heightDiff, horizontalDistance)
	tiltAngleDeg := tiltAngleRad * 180 / math.Pi

	fmt.Printf("\nStep 6 - Calculate Tilt Angle:\n")
	fmt.Printf("  Horizontal Distance: %.2f meters\n", horizontalDistance)
	fmt.Printf("  Height Difference: %.2f meters (camera %.0fm above target)\n", heightDiff, cameraHeight)
	fmt.Printf("  Tilt Angle: %.2f degrees (%.4f radians)\n", tiltAngleDeg, tiltAngleRad)
	if tiltAngleDeg < 0 {
		fmt.Printf("  Direction: Looking DOWN %.2f degrees from horizontal\n", -tiltAngleDeg)
	} else {
		fmt.Printf("  Direction: Looking UP %.2f degrees from horizontal\n", tiltAngleDeg)
	}

	// Step 7: Convert angles to PTZ coordinate system [-2, 0]
	targetX = -0.5 - (panAngleRad / math.Pi)
	targetY = -0.5 + (tiltAngleRad / (math.Pi / 2))

	fmt.Printf("\nStep 7 - Convert to PTZ Coordinates (before wrapping):\n")
	fmt.Printf("  Target X (Pan):  %.4f\n", targetX)
	fmt.Printf("  Target Y (Tilt): %.4f\n", targetY)

	// Step 8: Wrap X to [-2, 0] range
	if targetX < -2.0 {
		targetX += 2.0
		fmt.Printf("  Wrapped X: %.4f (was < -2.0)\n", targetX)
	} else if targetX > 0.0 {
		targetX -= 2.0
		fmt.Printf("  Wrapped X: %.4f (was > 0.0)\n", targetX)
	}

	// Step 9: Wrap Y to [-2, 0] range
	if targetY < -2.0 {
		targetY += 2.0
		fmt.Printf("  Wrapped Y: %.4f (was < -2.0)\n", targetY)
	} else if targetY > 0.0 {
		targetY -= 2.0
		fmt.Printf("  Wrapped Y: %.4f (was > 0.0)\n", targetY)
	}

	fmt.Printf("\n========== Final PTZ Coordinates ==========\n")
	fmt.Printf("Pan (X):  %.4f ", targetX)
	if targetX == -0.5 {
		fmt.Printf("(Straight ahead)\n")
	} else if targetX > -0.5 && targetX < 0 {
		fmt.Printf("(Left of center)\n")
	} else if targetX < -0.5 && targetX > -1.0 {
		fmt.Printf("(Right of center)\n")
	} else if targetX == -1.0 {
		fmt.Printf("(90° right)\n")
	} else if targetX < -1.0 && targetX > -1.5 {
		fmt.Printf("(Right-rear)\n")
	} else if targetX == -1.5 {
		fmt.Printf("(Directly behind)\n")
	} else {
		fmt.Printf("(Left-rear)\n")
	}

	fmt.Printf("Tilt (Y): %.4f ", targetY)
	if targetY == -0.5 {
		fmt.Printf("(Level/Horizon)\n")
	} else if targetY > -0.5 && targetY < 0 {
		fmt.Printf("(Looking up)\n")
	} else if targetY < -0.5 && targetY > -1.0 {
		fmt.Printf("(Looking down)\n")
	} else if targetY == -1.0 {
		fmt.Printf("(Looking straight down)\n")
	} else {
		fmt.Printf("(Extreme position)\n")
	}
	fmt.Printf("Zoom:     %.1f\n", 0.0)
	fmt.Println("==========================================\n")

	return targetX, targetY
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║           GeoMove2 - Geographic PTZ Test                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Test case from your request
	testRequest := GeoMove2Request{
		ProfileToken:    "M1_Profile1",
		TargetLatitude:  30.72316,
		TargetLongitude: 76.73035,
		SelfLatitude:    30.72172,
		SelfLongitude:   76.73042,
		SelfHeading:     0,
		CameraHeight:    0.5,
		Zoom:            0.0,
	}

	fmt.Println("\n📍 Input Parameters:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Profile Token:      %s\n", testRequest.ProfileToken)
	fmt.Printf("\n🎯 Target Position:\n")
	fmt.Printf("  Latitude:         %.6f°\n", testRequest.TargetLatitude)
	fmt.Printf("  Longitude:        %.6f°\n", testRequest.TargetLongitude)
	fmt.Printf("\n📷 Camera Position:\n")
	fmt.Printf("  Latitude:         %.6f°\n", testRequest.SelfLatitude)
	fmt.Printf("  Longitude:        %.6f°\n", testRequest.SelfLongitude)
	fmt.Printf("  Heading:          %d° (0°=North, 90°=East, 180°=South, 270°=West)\n", testRequest.SelfHeading)
	fmt.Printf("  Camera Height:    %.1f meters above ground\n", testRequest.CameraHeight)
	fmt.Printf("  Zoom:             %.1f\n", testRequest.Zoom)

	// Calculate relative position
	latDiff := testRequest.TargetLatitude - testRequest.SelfLatitude
	longDiff := testRequest.TargetLongitude - testRequest.SelfLongitude

	fmt.Printf("\n📐 Quick Analysis:\n")
	if latDiff > 0 {
		fmt.Printf("  Target is %.1f meters NORTH\n", latDiff*metersPerDegreeLatitude())
	} else {
		fmt.Printf("  Target is %.1f meters SOUTH\n", -latDiff*metersPerDegreeLatitude())
	}

	if longDiff > 0 {
		fmt.Printf("  Target is %.1f meters EAST\n", longDiff*metersPerDegreeLongitude(testRequest.SelfLatitude))
	} else {
		fmt.Printf("  Target is %.1f meters WEST\n", -longDiff*metersPerDegreeLongitude(testRequest.SelfLatitude))
	}

	// Convert to PTZ coordinates
	targetX, targetY := convertLatLongToXY(
		testRequest.TargetLatitude,
		testRequest.TargetLongitude,
		testRequest.SelfHeading,
		testRequest.SelfLatitude,
		testRequest.SelfLongitude,
		testRequest.CameraHeight,
	)

	// Display as AbsoluteMove command
	fmt.Println("\n📤 Resulting AbsoluteMove Command:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("{\n")
	fmt.Printf("  \"profileToken\": \"%s\",\n", testRequest.ProfileToken)
	fmt.Printf("  \"position\": {\n")
	fmt.Printf("    \"panTilt\": {\n")
	fmt.Printf("      \"x\": %.6f,\n", targetX)
	fmt.Printf("      \"y\": %.6f\n", targetY)
	fmt.Printf("    },\n")
	fmt.Printf("    \"zoom\": {\n")
	fmt.Printf("      \"x\": %.1f\n", testRequest.Zoom)
	fmt.Printf("    }\n")
	fmt.Printf("  },\n")
	fmt.Printf("  \"speed\": {\n")
	fmt.Printf("    \"panTilt\": { \"x\": 0.5, \"y\": 0.5 },\n")
	fmt.Printf("    \"zoom\": { \"x\": 0.5 }\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")

	fmt.Println("\n✅ Test completed successfully!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}
