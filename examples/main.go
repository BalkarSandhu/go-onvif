package main

import (
	"fmt"
	"log"

	"github.com/use-go/onvif"
)

func main() {
	// Replace with your network interface name if needed
	interfaces := []string{"Ethernet 2"}

	// Discovery
	for _, iface := range interfaces {
		devices, err := onvif.GetAvailableDevicesAtSpecificEthernetInterface(iface)
		if err != nil {
			log.Printf("Failed on %s: %v\n", iface, err)
			continue
		}
		fmt.Printf("Devices found on %s: %v\n", iface, devices)
	}

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
		log.Fatalf("Failed to connect to device: %v", err)
	}

	fmt.Printf("Device Params: %v\n", dev.GetDeviceParams())
	fmt.Printf("Services: %v\n", dev.GetServices())
}
