package service

import (
	"context"
	"encoding/json"
	"fmt"
	onvif "go-onvif/internal"
	"go-onvif/internal/device"
	"go-onvif/internal/media"
	"go-onvif/internal/ptz"
	device_rpc "go-onvif/internal/sdk/device"
	media_rpc "go-onvif/internal/sdk/media"
	ptz_rpc "go-onvif/internal/sdk/ptz"
	xsd_onvif "go-onvif/internal/xsd/onvif"

	"math"
	"time"

	"github.com/juju/errors"
)

func CallDeviceMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodName {

	// ──────────────────── Basic / capability queries ────────────────────
	case "GetServices":
		var req device.GetServices
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetServices(ctx, dev, req)

	case "GetServiceCapabilities":
		return device_rpc.Call_GetServiceCapabilities(ctx, dev, device.GetServiceCapabilities{})

	case "GetDeviceInformation":
		return device_rpc.Call_GetDeviceInformation(ctx, dev, device.GetDeviceInformation{})

	// ──────────────────────── Date / time ───────────────────────────────
	case "SetSystemDateAndTime":
		var req device.SetSystemDateAndTime
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetSystemDateAndTime(ctx, dev, req)

	case "GetSystemDateAndTime":
		return device_rpc.Call_GetSystemDateAndTime(ctx, dev, device.GetSystemDateAndTime{})

	case "SetSystemFactoryDefault":
		var req device.SetSystemFactoryDefault
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetSystemFactoryDefault(ctx, dev, req)

	case "SystemReboot":
		return device_rpc.Call_SystemReboot(ctx, dev, device.SystemReboot{})

	// ───────────────────── System backup / log ──────────────────────────
	case "GetSystemBackup":
		var req device.GetSystemBackup
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetSystemBackup(ctx, dev, req)

	case "GetSystemLog":
		var req device.GetSystemLog
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetSystemLog(ctx, dev, req)

	case "GetSystemSupportInformation":
		return device_rpc.Call_GetSystemSupportInformation(ctx, dev, device.GetSystemSupportInformation{})

	// ───────────────────────── Scopes ───────────────────────────────────
	case "GetScopes":
		return device_rpc.Call_GetScopes(ctx, dev, device.GetScopes{})

	case "SetScopes":
		var req device.SetScopes
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetScopes(ctx, dev, req)

	case "AddScopes":
		var req device.AddScopes
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_AddScopes(ctx, dev, req)

	case "RemoveScopes":
		var req device.RemoveScopes
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_RemoveScopes(ctx, dev, req)

	// ───── Discovery modes & remote discovery ─────
	case "GetDiscoveryMode":
		return device_rpc.Call_GetDiscoveryMode(ctx, dev, device.GetDiscoveryMode{})

	case "SetDiscoveryMode":
		var req device.SetDiscoveryMode
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDiscoveryMode(ctx, dev, req)

	case "GetRemoteDiscoveryMode":
		return device_rpc.Call_GetRemoteDiscoveryMode(ctx, dev, device.GetRemoteDiscoveryMode{})

	case "SetRemoteDiscoveryMode":
		var req device.SetRemoteDiscoveryMode
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRemoteDiscoveryMode(ctx, dev, req)

	// ───────────────────── Remote user ──────────────────────────────────
	case "GetRemoteUser":
		return device_rpc.Call_GetRemoteUser(ctx, dev, device.GetRemoteUser{})

	case "SetRemoteUser":
		var req device.SetRemoteUser
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRemoteUser(ctx, dev, req)

	// ───────────────────────── Users ────────────────────────────────────
	case "GetUsers":
		return device_rpc.Call_GetUsers(ctx, dev, device.GetUsers{})

	case "CreateUsers":
		var req device.CreateUsers
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateUsers(ctx, dev, req)

	case "DeleteUsers":
		var req device.DeleteUsers
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteUsers(ctx, dev, req)

	case "SetUser":
		var req device.SetUser
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetUser(ctx, dev, req)

	// ───────────────────────── Capabilities / hostname ──────────────────
	case "GetWsdlUrl":
		return device_rpc.Call_GetWsdlUrl(ctx, dev, device.GetWsdlUrl{})

	case "GetCapabilities":
		var req device.GetCapabilities
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetCapabilities(ctx, dev, req)

	case "GetHostname":
		return device_rpc.Call_GetHostname(ctx, dev, device.GetHostname{})

	case "SetHostname":
		var req device.SetHostname
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetHostname(ctx, dev, req)

	case "SetHostnameFromDHCP":
		var req device.SetHostnameFromDHCP
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetHostnameFromDHCP(ctx, dev, req)

	// ───────────────────── DNS / NTP / DynDNS ───────────────────────────
	case "GetDNS":
		return device_rpc.Call_GetDNS(ctx, dev, device.GetDNS{})

	case "SetDNS":
		var req device.SetDNS
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDNS(ctx, dev, req)

	case "GetNTP":
		return device_rpc.Call_GetNTP(ctx, dev, device.GetNTP{})

	case "SetNTP":
		var req device.SetNTP
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNTP(ctx, dev, req)

	case "GetDynamicDNS":
		return device_rpc.Call_GetDynamicDNS(ctx, dev, device.GetDynamicDNS{})

	case "SetDynamicDNS":
		var req device.SetDynamicDNS
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDynamicDNS(ctx, dev, req)

	// ───────────────────── Network interfaces / protocols ───────────────
	case "GetNetworkInterfaces":
		return device_rpc.Call_GetNetworkInterfaces(ctx, dev, device.GetNetworkInterfaces{})

	case "SetNetworkInterfaces":
		var req device.SetNetworkInterfaces
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkInterfaces(ctx, dev, req)

	case "GetNetworkProtocols":
		return device_rpc.Call_GetNetworkProtocols(ctx, dev, device.GetNetworkProtocols{})

	case "SetNetworkProtocols":
		var req device.SetNetworkProtocols
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkProtocols(ctx, dev, req)

	case "GetNetworkDefaultGateway":
		return device_rpc.Call_GetNetworkDefaultGateway(ctx, dev, device.GetNetworkDefaultGateway{})

	case "SetNetworkDefaultGateway":
		var req device.SetNetworkDefaultGateway
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkDefaultGateway(ctx, dev, req)

	case "GetZeroConfiguration":
		return device_rpc.Call_GetZeroConfiguration(ctx, dev, device.GetZeroConfiguration{})

	case "SetZeroConfiguration":
		var req device.SetZeroConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetZeroConfiguration(ctx, dev, req)

	// ────────────────────── IP Address filter ───────────────────────────
	case "GetIPAddressFilter":
		return device_rpc.Call_GetIPAddressFilter(ctx, dev, device.GetIPAddressFilter{})

	case "SetIPAddressFilter":
		var req device.SetIPAddressFilter
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetIPAddressFilter(ctx, dev, req)

	case "AddIPAddressFilter":
		var req device.AddIPAddressFilter
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_AddIPAddressFilter(ctx, dev, req)

	case "RemoveIPAddressFilter":
		var req device.RemoveIPAddressFilter
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_RemoveIPAddressFilter(ctx, dev, req)

	// ─────────────────── Certificates / PKI  ────────────────────────────
	case "GetAccessPolicy":
		return device_rpc.Call_GetAccessPolicy(ctx, dev, device.GetAccessPolicy{})

	case "SetAccessPolicy":
		var req device.SetAccessPolicy
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetAccessPolicy(ctx, dev, req)

	case "CreateCertificate":
		var req device.CreateCertificate
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateCertificate(ctx, dev, req)

	case "GetCertificates":
		return device_rpc.Call_GetCertificates(ctx, dev, device.GetCertificates{})

	case "GetCertificatesStatus":
		return device_rpc.Call_GetCertificatesStatus(ctx, dev, device.GetCertificatesStatus{})

	case "SetCertificatesStatus":
		var req device.SetCertificatesStatus
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetCertificatesStatus(ctx, dev, req)

	case "DeleteCertificates":
		var req device.DeleteCertificates
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteCertificates(ctx, dev, req)

	case "GetPkcs10Request":
		var req device.GetPkcs10Request
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetPkcs10Request(ctx, dev, req)

	case "LoadCertificates":
		var req device.LoadCertificates
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCertificates(ctx, dev, req)

	case "GetClientCertificateMode":
		return device_rpc.Call_GetClientCertificateMode(ctx, dev, device.GetClientCertificateMode{})

	case "SetClientCertificateMode":
		var req device.SetClientCertificateMode
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetClientCertificateMode(ctx, dev, req)

	case "GetCACertificates":
		return device_rpc.Call_GetCACertificates(ctx, dev, device.GetCACertificates{})

	case "LoadCertificateWithPrivateKey":
		var req device.LoadCertificateWithPrivateKey
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCertificateWithPrivateKey(ctx, dev, req)

	case "GetCertificateInformation":
		var req device.GetCertificateInformation
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetCertificateInformation(ctx, dev, req)

	case "LoadCACertificates":
		var req device.LoadCACertificates
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCACertificates(ctx, dev, req)

	// ────────────────────── IEEE 802.1X / 802.11 ────────────────────────
	case "CreateDot1XConfiguration":
		var req device.CreateDot1XConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateDot1XConfiguration(ctx, dev, req)

	case "SetDot1XConfiguration":
		var req device.SetDot1XConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDot1XConfiguration(ctx, dev, req)

	case "GetDot1XConfiguration":
		var req device.GetDot1XConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetDot1XConfiguration(ctx, dev, req)

	case "GetDot1XConfigurations":
		return device_rpc.Call_GetDot1XConfigurations(ctx, dev, device.GetDot1XConfigurations{})

	case "DeleteDot1XConfiguration":
		var req device.DeleteDot1XConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteDot1XConfiguration(ctx, dev, req)

	case "GetDot11Capabilities":
		return device_rpc.Call_GetDot11Capabilities(ctx, dev, device.GetDot11Capabilities{})

	case "GetDot11Status":
		var req device.GetDot11Status
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetDot11Status(ctx, dev, req)

	case "ScanAvailableDot11Networks":
		var req device.ScanAvailableDot11Networks
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_ScanAvailableDot11Networks(ctx, dev, req)

	// ───────────────────── Relay / auxiliary ────────────────────────────
	case "GetRelayOutputs":
		return device_rpc.Call_GetRelayOutputs(ctx, dev, device.GetRelayOutputs{})

	case "SetRelayOutputSettings":
		var req device.SetRelayOutputSettings
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRelayOutputSettings(ctx, dev, req)

	case "SetRelayOutputState":
		var req device.SetRelayOutputState
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRelayOutputState(ctx, dev, req)

	case "SendAuxiliaryCommand":
		var req device.SendAuxiliaryCommand
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SendAuxiliaryCommand(ctx, dev, req)

	// ───────────────────── Firmware / restore ───────────────────────────
	case "GetSystemUris":
		return device_rpc.Call_GetSystemUris(ctx, dev, device.GetSystemUris{})

	case "StartFirmwareUpgrade":
		var req device.StartFirmwareUpgrade
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_StartFirmwareUpgrade(ctx, dev, req)

	case "StartSystemRestore":
		var req device.StartSystemRestore
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_StartSystemRestore(ctx, dev, req)

	// ───────────────────── Storage configuration ────────────────────────
	case "GetStorageConfigurations":
		return device_rpc.Call_GetStorageConfigurations(ctx, dev, device.GetStorageConfigurations{})

	case "CreateStorageConfiguration":
		var req device.CreateStorageConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateStorageConfiguration(ctx, dev, req)

	case "GetStorageConfiguration":
		var req device.GetStorageConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetStorageConfiguration(ctx, dev, req)

	case "SetStorageConfiguration":
		var req device.SetStorageConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetStorageConfiguration(ctx, dev, req)

	case "DeleteStorageConfiguration":
		var req device.DeleteStorageConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteStorageConfiguration(ctx, dev, req)

	// ───────────────────── Geo-location (ONVIF Profile G) ───────────────
	case "GetGeoLocation":
		return device_rpc.Call_GetGeoLocation(ctx, dev, device.GetGeoLocation{})

	case "SetGeoLocation":
		var req device.SetGeoLocation
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetGeoLocation(ctx, dev, req)

	case "DeleteGeoLocation":
		var req device.DeleteGeoLocation
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteGeoLocation(ctx, dev, req)

	// ─────────────────────────────────────────────────────────────────────
	default:
		return nil, errors.New("unknown device method: " + methodName)
	}
}

func CallPTZMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodName {
	case "GetServiceCapabilities":
		return ptz_rpc.Call_GetServiceCapabilities(ctx, dev, ptz.GetServiceCapabilities{})

	case "GetNodes":
		return ptz_rpc.Call_GetNodes(ctx, dev, ptz.GetNodes{})

	case "GetNode":
		var request ptz.GetNode
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetNode(ctx, dev, request)

	case "GetConfiguration":
		var request ptz.GetConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetConfiguration(ctx, dev, request)

	case "GetConfigurations":
		return ptz_rpc.Call_GetConfigurations(ctx, dev, ptz.GetConfigurations{})

	case "SetConfiguration":
		var request ptz.SetConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_SetConfiguration(ctx, dev, request)

	case "GetConfigurationOptions":
		var request ptz.GetConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetConfigurationOptions(ctx, dev, request)

	case "SendAuxiliaryCommand":
		var request ptz.SendAuxiliaryCommand
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_SendAuxiliaryCommand(ctx, dev, request)

	case "GetPresets":
		var request ptz.GetPresets
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetPresets(ctx, dev, request)

	case "SetPreset":
		var request ptz.SetPreset
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_SetPreset(ctx, dev, request)

	case "RemovePreset":
		var request ptz.RemovePreset
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_RemovePreset(ctx, dev, request)

	case "GotoPreset":
		var request ptz.GotoPreset
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GotoPreset(ctx, dev, request)

	case "GotoHomePosition":
		var request ptz.GotoHomePosition
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GotoHomePosition(ctx, dev, request)

	case "SetHomePosition":
		var request ptz.SetHomePosition
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_SetHomePosition(ctx, dev, request)

	case "ContinuousMove":
		var request ptz.ContinuousMove
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_ContinuousMove(ctx, dev, request)

	case "RelativeMove":
		var request ptz.RelativeMove
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_RelativeMove(ctx, dev, request)

	case "GetStatus":
		var request ptz.GetStatus
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetStatus(ctx, dev, request)

	case "AbsoluteMove":
		var request ptz.AbsoluteMove
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_AbsoluteMove(ctx, dev, request)

	case "GeoMove":
		var request ptz.GeoMove
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GeoMove(ctx, dev, request)

	case "Stop":
		var request ptz.Stop
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_Stop(ctx, dev, request)

	case "GetPresetTours":
		return ptz_rpc.Call_GetPresetTours(ctx, dev, ptz.GetPresetTours{})

	case "GetPresetTour":
		var request ptz.GetPresetTour
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetPresetTour(ctx, dev, request)

	case "GetPresetTourOptions":
		var request ptz.GetPresetTourOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetPresetTourOptions(ctx, dev, request)

	case "CreatePresetTour":
		var request ptz.CreatePresetTour
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_CreatePresetTour(ctx, dev, request)

	case "ModifyPresetTour":
		var request ptz.ModifyPresetTour
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_ModifyPresetTour(ctx, dev, request)

	case "OperatePresetTour":
		var request ptz.OperatePresetTour
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_OperatePresetTour(ctx, dev, request)

	case "RemovePresetTour":
		var request ptz.RemovePresetTour
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_RemovePresetTour(ctx, dev, request)

	case "GetCompatibleConfigurations":
		var request ptz.GetCompatibleConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetCompatibleConfigurations(ctx, dev, request)
	case "AbsoluteMove2":
		var request AbsoluteMove2
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}

		// Extract target position from request
		targetPan := request.Position.PanTilt.X
		targetTilt := request.Position.PanTilt.Y
		targetZoom := request.Position.Zoom.X

		// Configuration
		const (
			POSITION_TOLERANCE = 0.05 // Acceptable position error
			MAX_SPEED          = 0.5  // Maximum movement speed
			MIN_SPEED          = 0.1  // Minimum speed to ensure movement
			TIMEOUT            = 30 * time.Second
			POLL_INTERVAL      = 100 * time.Millisecond
		)

		startTime := time.Now()

		// Position control loop
		for {
			// 1. Check timeout
			if time.Since(startTime) > TIMEOUT {
				// Stop camera
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return nil, fmt.Errorf("timeout: could not reach target position")
			}

			// 2. Get current position
			statusReq := ptz.GetStatus{
				ProfileToken: request.ProfileToken,
			}
			status, err := ptz_rpc.Call_GetStatus(ctx, dev, statusReq)
			if err != nil {
				return nil, fmt.Errorf("failed to get status: %w", err)
			}

			currentPan := status.PTZStatus.Position.PanTilt.X
			currentTilt := status.PTZStatus.Position.PanTilt.Y
			currentZoom := status.PTZStatus.Position.Zoom.X

			// 3. Calculate distances
			panDistance := targetPan - currentPan
			tiltDistance := targetTilt - currentTilt
			zoomDistance := targetZoom - currentZoom

			// 4. Check if target reached
			panReached := math.Abs(panDistance) < POSITION_TOLERANCE
			tiltReached := math.Abs(tiltDistance) < POSITION_TOLERANCE
			zoomReached := math.Abs(zoomDistance) < POSITION_TOLERANCE

			if panReached && tiltReached && zoomReached {
				// Stop camera
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return status, nil // Success - return final status
			}

			// 5. Calculate speeds (proportional control)
			calculateSpeed := func(distance float64) float64 {
				if math.Abs(distance) < POSITION_TOLERANCE {
					return 0.0
				}

				// Proportional speed based on distance
				speed := distance * 2.0 // Gain factor

				// Clamp to max speed
				if math.Abs(speed) > MAX_SPEED {
					if speed > 0 {
						speed = MAX_SPEED
					} else {
						speed = -MAX_SPEED
					}
				} else if math.Abs(speed) < MIN_SPEED {
					// Ensure minimum speed if not at target
					if distance > 0 {
						speed = MIN_SPEED
					} else {
						speed = -MIN_SPEED
					}
				}

				return speed
			}

			panSpeed := calculateSpeed(panDistance)
			tiltSpeed := calculateSpeed(tiltDistance)
			zoomSpeed := calculateSpeed(zoomDistance)

			// 6. Send continuous move command
			moveReq := ptz.ContinuousMove{
				ProfileToken: request.ProfileToken,
				Velocity: xsd_onvif.PTZSpeed{
					PanTilt: xsd_onvif.Vector2D{
						X: panSpeed,
						Y: tiltSpeed,
					},
					Zoom: xsd_onvif.Vector1D{
						X: zoomSpeed,
					},
				},
			}

			if _, err := ptz_rpc.Call_ContinuousMove(ctx, dev, moveReq); err != nil {
				// Stop on error
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return nil, fmt.Errorf("continuous move failed: %w", err)
			}

			// 7. Wait before next iteration
			time.Sleep(POLL_INTERVAL)
		}

	case "GeoMove2":
		var request GeoMove2
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}

		target := convertLatLongToXY(
			request.TargetLatitude,
			request.TargetLongitude,
			request.SelfHeading,
			request.SelfLatitude,
			request.SelfLongitude,
			request.CameraHeight,
		)

		targetPan := target.X
		targetTilt := target.Y
		targetZoom := request.Zoom

		// Configuration
		const (
			POSITION_TOLERANCE = 0.1 // Acceptable position error
			MAX_SPEED          = 0.3 // Maximum movement speed
			MIN_SPEED          = 0.1 // Minimum speed to ensure movement
			TIMEOUT            = 30 * time.Second
			POLL_INTERVAL      = 100 * time.Millisecond
		)

		startTime := time.Now()

		// Position control loop
		for {
			// 1. Check timeout
			if time.Since(startTime) > TIMEOUT {
				// Stop camera
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return nil, fmt.Errorf("timeout: could not reach target position")
			}

			// 2. Get current position
			statusReq := ptz.GetStatus{
				ProfileToken: request.ProfileToken,
			}
			status, err := ptz_rpc.Call_GetStatus(ctx, dev, statusReq)
			if err != nil {
				return nil, fmt.Errorf("failed to get status: %w", err)
			}

			currentPan := status.PTZStatus.Position.PanTilt.X
			currentTilt := status.PTZStatus.Position.PanTilt.Y
			currentZoom := status.PTZStatus.Position.Zoom.X

			// 3. Calculate distances
			panDistance := targetPan - currentPan
			tiltDistance := targetTilt - currentTilt
			zoomDistance := targetZoom - currentZoom

			// 4. Check if target reached
			panReached := math.Abs(panDistance) < POSITION_TOLERANCE
			tiltReached := math.Abs(tiltDistance) < POSITION_TOLERANCE
			zoomReached := math.Abs(zoomDistance) < POSITION_TOLERANCE

			if panReached && tiltReached && zoomReached {
				// Stop camera
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return status, nil // Success - return final status
			}

			// 5. Calculate speeds (proportional control)
			calculateSpeed := func(distance float64) float64 {
				if math.Abs(distance) < POSITION_TOLERANCE {
					return 0.0
				}

				// Proportional speed based on distance
				speed := distance * 2.0 // Gain factor

				// Clamp to max speed
				if math.Abs(speed) > MAX_SPEED {
					if speed > 0 {
						speed = MAX_SPEED
					} else {
						speed = -MAX_SPEED
					}
				} else if math.Abs(speed) < MIN_SPEED {
					// Ensure minimum speed if not at target
					if distance > 0 {
						speed = MIN_SPEED
					} else {
						speed = -MIN_SPEED
					}
				}

				return speed
			}

			panSpeed := calculateSpeed(panDistance)
			tiltSpeed := calculateSpeed(tiltDistance)
			zoomSpeed := calculateSpeed(zoomDistance)

			// 6. Send continuous move command
			moveReq := ptz.ContinuousMove{
				ProfileToken: request.ProfileToken,
				Velocity: xsd_onvif.PTZSpeed{
					PanTilt: xsd_onvif.Vector2D{
						X: panSpeed,
						Y: tiltSpeed,
					},
					Zoom: xsd_onvif.Vector1D{
						X: zoomSpeed,
					},
				},
			}

			if _, err := ptz_rpc.Call_ContinuousMove(ctx, dev, moveReq); err != nil {
				// Stop on error
				stopReq := ptz.Stop{
					ProfileToken: request.ProfileToken,
					PanTilt:      true,
					Zoom:         true,
				}
				ptz_rpc.Call_Stop(ctx, dev, stopReq)
				return nil, fmt.Errorf("continuous move failed: %w", err)
			}

			// 7. Wait before next iteration
			time.Sleep(POLL_INTERVAL)
		}

	default:
		return nil, errors.New("unknown PTZ method: " + methodName)
	}
}

// callMediaMethod routes a Media service method name to the corresponding RPC
// call. If the RPC needs a request payload, it is unmarshalled from `data`.
func CallMediaMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodName {

	// ---------------------------------------------------------------------
	// ――― Simple “getter” methods that take no parameters ―――
	// ---------------------------------------------------------------------
	case "GetServiceCapabilities":
		return media_rpc.Call_GetServiceCapabilities(ctx, dev, media.GetServiceCapabilities{})
	case "GetVideoSources":
		return media_rpc.Call_GetVideoSources(ctx, dev, media.GetVideoSources{})
	case "GetAudioSources":
		return media_rpc.Call_GetAudioSources(ctx, dev, media.GetAudioSources{})
	case "GetAudioOutputs":
		return media_rpc.Call_GetAudioOutputs(ctx, dev, media.GetAudioOutputs{})
	case "GetProfile":
		return media_rpc.Call_GetProfile(ctx, dev, media.GetProfile{})
	case "GetProfiles":
		return media_rpc.Call_GetProfiles(ctx, dev, media.GetProfiles{})
	case "GetVideoSourceConfigurations":
		return media_rpc.Call_GetVideoSourceConfigurations(ctx, dev, media.GetVideoSourceConfigurations{})
	case "GetVideoEncoderConfigurations":
		return media_rpc.Call_GetVideoEncoderConfigurations(ctx, dev, media.GetVideoEncoderConfigurations{})
	case "GetAudioSourceConfigurations":
		return media_rpc.Call_GetAudioSourceConfigurations(ctx, dev, media.GetAudioSourceConfigurations{})
	case "GetAudioEncoderConfigurations":
		return media_rpc.Call_GetAudioEncoderConfigurations(ctx, dev, media.GetAudioEncoderConfigurations{})
	case "GetVideoAnalyticsConfigurations":
		return media_rpc.Call_GetVideoAnalyticsConfigurations(ctx, dev, media.GetVideoAnalyticsConfigurations{})
	case "GetMetadataConfigurations":
		return media_rpc.Call_GetMetadataConfigurations(ctx, dev, media.GetMetadataConfigurations{})
	case "GetAudioOutputConfigurations":
		return media_rpc.Call_GetAudioOutputConfigurations(ctx, dev, media.GetAudioOutputConfigurations{})
	case "GetAudioDecoderConfigurations":
		return media_rpc.Call_GetAudioDecoderConfigurations(ctx, dev, media.GetAudioDecoderConfigurations{})
	case "GetVideoSourceConfiguration":
		return media_rpc.Call_GetVideoSourceConfiguration(ctx, dev, media.GetVideoSourceConfiguration{})
	case "GetVideoEncoderConfiguration":
		return media_rpc.Call_GetVideoEncoderConfiguration(ctx, dev, media.GetVideoEncoderConfiguration{})
	case "GetAudioSourceConfiguration":
		return media_rpc.Call_GetAudioSourceConfiguration(ctx, dev, media.GetAudioSourceConfiguration{})
	case "GetAudioEncoderConfiguration":
		return media_rpc.Call_GetAudioEncoderConfiguration(ctx, dev, media.GetAudioEncoderConfiguration{})
	case "GetVideoAnalyticsConfiguration":
		return media_rpc.Call_GetVideoAnalyticsConfiguration(ctx, dev, media.GetVideoAnalyticsConfiguration{})
	case "GetMetadataConfiguration":
		return media_rpc.Call_GetMetadataConfiguration(ctx, dev, media.GetMetadataConfiguration{})
	case "GetAudioOutputConfiguration":
		return media_rpc.Call_GetAudioOutputConfiguration(ctx, dev, media.GetAudioOutputConfiguration{})
	case "GetAudioDecoderConfiguration":
		return media_rpc.Call_GetAudioDecoderConfiguration(ctx, dev, media.GetAudioDecoderConfiguration{})
	case "GetCompatibleVideoEncoderConfigurations":
		return media_rpc.Call_GetCompatibleVideoEncoderConfigurations(ctx, dev, media.GetCompatibleVideoEncoderConfigurations{})
	case "GetCompatibleVideoSourceConfigurations":
		return media_rpc.Call_GetCompatibleVideoSourceConfigurations(ctx, dev, media.GetCompatibleVideoSourceConfigurations{})
	case "GetCompatibleAudioEncoderConfigurations":
		return media_rpc.Call_GetCompatibleAudioEncoderConfigurations(ctx, dev, media.GetCompatibleAudioEncoderConfigurations{})
	case "GetCompatibleAudioSourceConfigurations":
		return media_rpc.Call_GetCompatibleAudioSourceConfigurations(ctx, dev, media.GetCompatibleAudioSourceConfigurations{})
	case "GetCompatibleVideoAnalyticsConfigurations":
		return media_rpc.Call_GetCompatibleVideoAnalyticsConfigurations(ctx, dev, media.GetCompatibleVideoAnalyticsConfigurations{})
	case "GetCompatibleMetadataConfigurations":
		return media_rpc.Call_GetCompatibleMetadataConfigurations(ctx, dev, media.GetCompatibleMetadataConfigurations{})
	case "GetCompatibleAudioOutputConfigurations":
		return media_rpc.Call_GetCompatibleAudioOutputConfigurations(ctx, dev, media.GetCompatibleAudioOutputConfigurations{})
	case "GetCompatibleAudioDecoderConfigurations":
		return media_rpc.Call_GetCompatibleAudioDecoderConfigurations(ctx, dev, media.GetCompatibleAudioDecoderConfigurations{})
	case "GetGuaranteedNumberOfVideoEncoderInstances":
		return media_rpc.Call_GetGuaranteedNumberOfVideoEncoderInstances(ctx, dev, media.GetGuaranteedNumberOfVideoEncoderInstances{})

	// ---------------------------------------------------------------------
	// ――― Methods that need a payload (unmarshal from `data`) ―――
	// ---------------------------------------------------------------------
	case "CreateProfile":
		var req media.CreateProfile
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_CreateProfile(ctx, dev, req)

	case "AddVideoEncoderConfiguration":
		var req media.AddVideoEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoEncoderConfiguration(ctx, dev, req)
	case "RemoveVideoEncoderConfiguration":
		var req media.RemoveVideoEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoEncoderConfiguration(ctx, dev, req)

	case "AddVideoSourceConfiguration":
		var req media.AddVideoSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoSourceConfiguration(ctx, dev, req)
	case "RemoveVideoSourceConfiguration":
		var req media.RemoveVideoSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoSourceConfiguration(ctx, dev, req)

	case "AddAudioEncoderConfiguration":
		var req media.AddAudioEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioEncoderConfiguration(ctx, dev, req)
	case "RemoveAudioEncoderConfiguration":
		var req media.RemoveAudioEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioEncoderConfiguration(ctx, dev, req)

	case "AddAudioSourceConfiguration":
		var req media.AddAudioSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioSourceConfiguration(ctx, dev, req)
	case "RemoveAudioSourceConfiguration":
		var req media.RemoveAudioSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioSourceConfiguration(ctx, dev, req)

	case "AddPTZConfiguration":
		var req media.AddPTZConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddPTZConfiguration(ctx, dev, req)
	case "RemovePTZConfiguration":
		var req media.RemovePTZConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemovePTZConfiguration(ctx, dev, req)

	case "AddVideoAnalyticsConfiguration":
		var req media.AddVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoAnalyticsConfiguration(ctx, dev, req)
	case "RemoveVideoAnalyticsConfiguration":
		var req media.RemoveVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoAnalyticsConfiguration(ctx, dev, req)

	case "AddMetadataConfiguration":
		var req media.AddMetadataConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddMetadataConfiguration(ctx, dev, req)
	case "RemoveMetadataConfiguration":
		var req media.RemoveMetadataConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveMetadataConfiguration(ctx, dev, req)

	case "AddAudioOutputConfiguration":
		var req media.AddAudioOutputConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioOutputConfiguration(ctx, dev, req)
	case "RemoveAudioOutputConfiguration":
		var req media.RemoveAudioOutputConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioOutputConfiguration(ctx, dev, req)

	case "AddAudioDecoderConfiguration":
		var req media.AddAudioDecoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioDecoderConfiguration(ctx, dev, req)
	case "RemoveAudioDecoderConfiguration":
		var req media.RemoveAudioDecoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioDecoderConfiguration(ctx, dev, req)

	case "DeleteProfile":
		var req media.DeleteProfile
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_DeleteProfile(ctx, dev, req)

	// ---- “Set*” calls ----------------------------------------------------
	case "SetVideoSourceConfiguration":
		var req media.SetVideoSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoSourceConfiguration(ctx, dev, req)
	case "SetVideoEncoderConfiguration":
		var req media.SetVideoEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoEncoderConfiguration(ctx, dev, req)
	case "SetAudioSourceConfiguration":
		var req media.SetAudioSourceConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioSourceConfiguration(ctx, dev, req)
	case "SetAudioEncoderConfiguration":
		var req media.SetAudioEncoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioEncoderConfiguration(ctx, dev, req)
	case "SetVideoAnalyticsConfiguration":
		var req media.SetVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoAnalyticsConfiguration(ctx, dev, req)
	case "SetMetadataConfiguration":
		var req media.SetMetadataConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetMetadataConfiguration(ctx, dev, req)
	case "SetAudioOutputConfiguration":
		var req media.SetAudioOutputConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioOutputConfiguration(ctx, dev, req)
	case "SetAudioDecoderConfiguration":
		var req media.SetAudioDecoderConfiguration
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioDecoderConfiguration(ctx, dev, req)

	// ---- “Options” calls -------------------------------------------------
	case "GetVideoSourceConfigurationOptions":
		var req media.GetVideoSourceConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoSourceConfigurationOptions(ctx, dev, req)
	case "GetVideoEncoderConfigurationOptions":
		var req media.GetVideoEncoderConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoEncoderConfigurationOptions(ctx, dev, req)
	case "GetAudioSourceConfigurationOptions":
		var req media.GetAudioSourceConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioSourceConfigurationOptions(ctx, dev, req)
	case "GetAudioEncoderConfigurationOptions":
		var req media.GetAudioEncoderConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioEncoderConfigurationOptions(ctx, dev, req)
	case "GetMetadataConfigurationOptions":
		var req media.GetMetadataConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetMetadataConfigurationOptions(ctx, dev, req)
	case "GetAudioOutputConfigurationOptions":
		var req media.GetAudioOutputConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioOutputConfigurationOptions(ctx, dev, req)
	case "GetAudioDecoderConfigurationOptions":
		var req media.GetAudioDecoderConfigurationOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioDecoderConfigurationOptions(ctx, dev, req)

	// ---- Streaming helpers ----------------------------------------------
	case "GetStreamUri":
		var req media.GetStreamUri
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetStreamUri(ctx, dev, req)
	case "StartMulticastStreaming":
		var req media.StartMulticastStreaming
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_StartMulticastStreaming(ctx, dev, req)
	case "StopMulticastStreaming":
		var req media.StopMulticastStreaming
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_StopMulticastStreaming(ctx, dev, req)
	case "SetSynchronizationPoint":
		var req media.SetSynchronizationPoint
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetSynchronizationPoint(ctx, dev, req)
	case "GetSnapshotUri":
		var req media.GetSnapshotUri
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetSnapshotUri(ctx, dev, req)

	// ---- OSD (On-screen display) ----------------------------------------
	case "GetVideoSourceModes":
		var req media.GetVideoSourceModes
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoSourceModes(ctx, dev, req)
	case "SetVideoSourceMode":
		var req media.SetVideoSourceMode
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoSourceMode(ctx, dev, req)
	case "GetOSDs":
		var req media.GetOSDs
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSDs(ctx, dev, req)
	case "GetOSD":
		var req media.GetOSD
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSD(ctx, dev, req)
	case "GetOSDOptions":
		var req media.GetOSDOptions
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSDOptions(ctx, dev, req)
	case "SetOSD":
		var req media.SetOSD
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetOSD(ctx, dev, req)
	case "CreateOSD":
		var req media.CreateOSD
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_CreateOSD(ctx, dev, req)
	case "DeleteOSD":
		var req media.DeleteOSD
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return media_rpc.Call_DeleteOSD(ctx, dev, req)

	//----------------------------------------------------------------------
	default:
		return nil, errors.New("unknown Media method: " + methodName)
	}
}
