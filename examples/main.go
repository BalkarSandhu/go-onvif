package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/juju/errors"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/sdk"
)

// Call_GetCapabilities forwards the call to dev.CallMethod() then parses the payload of the reply as a GetCapabilitiesResponse.
func Call_GetCapabilities(ctx context.Context, dev *onvif.Device, request device.GetCapabilities) (device.GetCapabilitiesResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			GetCapabilitiesResponse device.GetCapabilitiesResponse
		}
	}
	var reply Envelope
	if httpReply, err := dev.CallMethod(request); err != nil {
		return reply.Body.GetCapabilitiesResponse, errors.Annotate(err, "call")
	} else {
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "GetCapabilities")
		return reply.Body.GetCapabilitiesResponse, errors.Annotate(err, "reply")
	}
}

func main() {

	// ONVIF Camera Details
	cameraIP := "192.168.29.109" // Update with your camera's IP
	username := "admin"          // Update with actual username
	password := "admin123"       // Update with actual password

	//login
	// Connect to ONVIF camera
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    cameraIP,
		Username: username,
		Password: password,
	})

	if err != nil {
		fmt.Printf("Failed to create device: %v\n", err)
	}

	// Get Capabilities
	resp, err := Call_GetCapabilities(context.Background(), dev, device.GetCapabilities{Category: "All"})
	if err != nil {
		fmt.Printf("Failed to get capabilities: %v\n", err)
	}

	fmt.Println(resp.Capabilities)

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Print the JSON response
	fmt.Println(string(jsonData))
}
