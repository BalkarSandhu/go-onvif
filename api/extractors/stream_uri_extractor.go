package extractors

import (
	"encoding/json"
	"fmt"
	"go-onvif/api/service"
	onvif "go-onvif/internal"
)

type StreamUri struct {
	Uri string `json:"Uri"`
}

type StreamInfo struct {
	MediaUri StreamUri `json:"MediaUri"`
}

func GetStreamURI(dev *onvif.Device, profileToken string) (string, error) {
	// Prepare JSON payload
	payload := map[string]string{
		"ProfileToken": profileToken,
	}
	profileTokenBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile token: %w", err)
	}

	// Call GetStreamUri method
	resp, err := service.CallMediaMethod("GetStreamUri", dev, profileTokenBytes)
	if err != nil {
		return "", fmt.Errorf("CallMediaMethod failed: %w", err)
	}

	// Convert response to proper structure
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	var streamInfo StreamInfo
	err = json.Unmarshal(respBytes, &streamInfo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal stream info: %w", err)
	}

	return streamInfo.MediaUri.Uri, nil
}
