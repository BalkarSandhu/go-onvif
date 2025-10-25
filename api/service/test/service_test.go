package service_test

import (
	"testing"
	"time"

	"go-onvif/api/service"
	"go-onvif/internal/ptz"
	"go-onvif/internal/xsd/onvif"

	"github.com/stretchr/testify/assert"
)

// Mock GetStatusResponse for testing
func mockStatus(pan, tilt, zoom float64) ptz.GetStatusResponse {
	return ptz.GetStatusResponse{
		PTZStatus:	onvif.PTZStatus{
			Position: onvif.PTZVector{
				PanTilt: onvif.Vector2D{X: pan, Y: tilt},
				Zoom:    onvif.Vector1D{X: zoom},
			},
		},
	}
}

func TestGeoMoveDeterministic(t *testing.T) {
	requestJSON := `{
		"ProfileToken": "M1_Profile1",
		"camera": {
			"latitude": 23.72393,
			"longitude": 86.75771,
			"height": 10.0,
			"yaw_offset": 0.0,
			"tilt_offset": 0.0,
			"pan_range": {"min": -180.0, "max": 180.0},
			"tilt_range": {"min": -90.0, "max": 90.0}
		},
		"target": {
			"latitude": 23.72380,
			"longitude": 86.75780,
			"height": 0.0
		},
		"settings": {
			"velocity_scale": 0.5,
			"zoom_velocity_scale": 0.3,
			"max_velocity": 1.0,
			"timeout": "PT0.5S",
			"min_zoom_distance": 10.0,
			"max_zoom_distance": 500.0,
			"deadband_pan": 0.5,
			"deadband_tilt": 0.3,
			"deadband_zoom": 0.02,
			"precision": 6
		}
	}`

	var velocities []struct {
		Pan  float64
		Tilt float64
		Zoom float64
	}

	// Simulate multiple iterations
	for i := 0; i < 5; i++ {
		status := mockStatus(0.0, 0.0, 0.0) // start at origin

		cmd, err := service.GetMoveParams(status, []byte(requestJSON))
		assert.NoError(t, err)

		velocities = append(velocities, struct {
			Pan  float64
			Tilt float64
			Zoom float64
		}{
			Pan:  cmd.Velocity.PanTilt.X,
			Tilt: cmd.Velocity.PanTilt.Y,
			Zoom: cmd.Velocity.Zoom.X,
		})
		time.Sleep(10 * time.Millisecond) // simulate loop delay
	}

	// Assert all velocity outputs are identical
	first := velocities[0]
	for i, v := range velocities {
		assert.Equal(t, first.Pan, v.Pan, "Pan mismatch at iteration %d", i)
		assert.Equal(t, first.Tilt, v.Tilt, "Tilt mismatch at iteration %d", i)
		assert.Equal(t, first.Zoom, v.Zoom, "Zoom mismatch at iteration %d", i)
	}
}
