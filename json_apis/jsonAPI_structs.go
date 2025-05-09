package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/BalkarSandhu/go-onvif/device"
	"github.com/BalkarSandhu/go-onvif/media"
	"github.com/BalkarSandhu/go-onvif/onvif"
	"github.com/BalkarSandhu/go-onvif/ptz"
	device_rpc "github.com/BalkarSandhu/go-onvif/sdk/device"
	media_rpc "github.com/BalkarSandhu/go-onvif/sdk/media"
	ptz_rpc "github.com/BalkarSandhu/go-onvif/sdk/ptz"
	"github.com/juju/errors"
)

// callDeviceMethod handles all device service methods
func callDeviceMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodName {
	case "GetServices":
		var request device.GetServices
		if err := json.Unmarshal(data, &request); err != nil {
			// If no data provided, use empty struct
			request = device.GetServices{}
		}
		return device_rpc.Call_GetServices(ctx, dev, request)
	case "GetServiceCapabilities":
		return device_rpc.Call_GetServiceCapabilities(ctx, dev, device.GetServiceCapabilities{})
	case "GetDeviceInformation":
		return device_rpc.Call_GetDeviceInformation(ctx, dev, device.GetDeviceInformation{})
	case "SetSystemDateAndTime":
		var request device.SetSystemDateAndTime
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetSystemDateAndTime(ctx, dev, request)
	case "GetSystemDateAndTime":
		return device_rpc.Call_GetSystemDateAndTime(ctx, dev, device.GetSystemDateAndTime{})
	case "SetSystemFactoryDefault":
		var request device.SetSystemFactoryDefault
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetSystemFactoryDefault(ctx, dev, request)
	case "SystemReboot":
		return device_rpc.Call_SystemReboot(ctx, dev, device.SystemReboot{})
	case "GetSystemBackup":
		return device_rpc.Call_GetSystemBackup(ctx, dev, device.GetSystemBackup{})
	case "GetSystemLog":
		var request device.GetSystemLog
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetSystemLog(ctx, dev, request)
	case "GetSystemSupportInformation":
		return device_rpc.Call_GetSystemSupportInformation(ctx, dev, device.GetSystemSupportInformation{})
	case "GetScopes":
		return device_rpc.Call_GetScopes(ctx, dev, device.GetScopes{})
	case "SetScopes":
		var request device.SetScopes
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetScopes(ctx, dev, request)
	case "AddScopes":
		var request device.AddScopes
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_AddScopes(ctx, dev, request)
	case "RemoveScopes":
		var request device.RemoveScopes
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_RemoveScopes(ctx, dev, request)
	case "GetDiscoveryMode":
		return device_rpc.Call_GetDiscoveryMode(ctx, dev, device.GetDiscoveryMode{})
	case "SetDiscoveryMode":
		var request device.SetDiscoveryMode
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDiscoveryMode(ctx, dev, request)
	case "GetRemoteDiscoveryMode":
		return device_rpc.Call_GetRemoteDiscoveryMode(ctx, dev, device.GetRemoteDiscoveryMode{})
	case "SetRemoteDiscoveryMode":
		var request device.SetRemoteDiscoveryMode
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRemoteDiscoveryMode(ctx, dev, request)
	case "GetDPAddresses":
		return device_rpc.Call_GetDPAddresses(ctx, dev, device.GetDPAddresses{})
	case "GetEndpointReference":
		return device_rpc.Call_GetEndpointReference(ctx, dev, device.GetEndpointReference{})
	case "GetRemoteUser":
		return device_rpc.Call_GetRemoteUser(ctx, dev, device.GetRemoteUser{})
	case "SetRemoteUser":
		var request device.SetRemoteUser
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRemoteUser(ctx, dev, request)
	case "GetUsers":
		return device_rpc.Call_GetUsers(ctx, dev, device.GetUsers{})
	case "CreateUsers":
		var request device.CreateUsers
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateUsers(ctx, dev, request)
	case "DeleteUsers":
		var request device.DeleteUsers
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteUsers(ctx, dev, request)
	case "SetUser":
		var request device.SetUser
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetUser(ctx, dev, request)
	case "GetWsdlUrl":
		return device_rpc.Call_GetWsdlUrl(ctx, dev, device.GetWsdlUrl{})
	case "GetCapabilities":
		var request device.GetCapabilities
		if err := json.Unmarshal(data, &request); err != nil {
			// If empty, use default
			request = device.GetCapabilities{}
		}
		return device_rpc.Call_GetCapabilities(ctx, dev, request)
	case "GetHostname":
		return device_rpc.Call_GetHostname(ctx, dev, device.GetHostname{})
	case "SetHostname":
		var request device.SetHostname
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetHostname(ctx, dev, request)
	case "SetHostnameFromDHCP":
		var request device.SetHostnameFromDHCP
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetHostnameFromDHCP(ctx, dev, request)
	case "GetDNS":
		return device_rpc.Call_GetDNS(ctx, dev, device.GetDNS{})
	case "SetDNS":
		var request device.SetDNS
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDNS(ctx, dev, request)
	case "GetNTP":
		return device_rpc.Call_GetNTP(ctx, dev, device.GetNTP{})
	case "SetNTP":
		var request device.SetNTP
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNTP(ctx, dev, request)
	case "GetDynamicDNS":
		return device_rpc.Call_GetDynamicDNS(ctx, dev, device.GetDynamicDNS{})
	case "SetDynamicDNS":
		var request device.SetDynamicDNS
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDynamicDNS(ctx, dev, request)
	case "GetNetworkInterfaces":
		return device_rpc.Call_GetNetworkInterfaces(ctx, dev, device.GetNetworkInterfaces{})
	case "SetNetworkInterfaces":
		var request device.SetNetworkInterfaces
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkInterfaces(ctx, dev, request)
	case "GetNetworkProtocols":
		return device_rpc.Call_GetNetworkProtocols(ctx, dev, device.GetNetworkProtocols{})
	case "SetNetworkProtocols":
		var request device.SetNetworkProtocols
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkProtocols(ctx, dev, request)
	case "GetNetworkDefaultGateway":
		return device_rpc.Call_GetNetworkDefaultGateway(ctx, dev, device.GetNetworkDefaultGateway{})
	case "SetNetworkDefaultGateway":
		var request device.SetNetworkDefaultGateway
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetNetworkDefaultGateway(ctx, dev, request)
	case "GetZeroConfiguration":
		return device_rpc.Call_GetZeroConfiguration(ctx, dev, device.GetZeroConfiguration{})
	case "SetZeroConfiguration":
		var request device.SetZeroConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetZeroConfiguration(ctx, dev, request)
	case "GetIPAddressFilter":
		return device_rpc.Call_GetIPAddressFilter(ctx, dev, device.GetIPAddressFilter{})
	case "SetIPAddressFilter":
		var request device.SetIPAddressFilter
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetIPAddressFilter(ctx, dev, request)
	case "AddIPAddressFilter":
		var request device.AddIPAddressFilter
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_AddIPAddressFilter(ctx, dev, request)
	case "RemoveIPAddressFilter":
		var request device.RemoveIPAddressFilter
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_RemoveIPAddressFilter(ctx, dev, request)
	case "GetAccessPolicy":
		return device_rpc.Call_GetAccessPolicy(ctx, dev, device.GetAccessPolicy{})
	case "SetAccessPolicy":
		var request device.SetAccessPolicy
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetAccessPolicy(ctx, dev, request)
	case "CreateCertificate":
		var request device.CreateCertificate
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateCertificate(ctx, dev, request)
	case "GetCertificates":
		return device_rpc.Call_GetCertificates(ctx, dev, device.GetCertificates{})
	case "GetCertificatesStatus":
		return device_rpc.Call_GetCertificatesStatus(ctx, dev, device.GetCertificatesStatus{})
	case "SetCertificatesStatus":
		var request device.SetCertificatesStatus
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetCertificatesStatus(ctx, dev, request)
	case "DeleteCertificates":
		var request device.DeleteCertificates
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteCertificates(ctx, dev, request)
	case "GetPkcs10Request":
		var request device.GetPkcs10Request
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetPkcs10Request(ctx, dev, request)
	case "LoadCertificates":
		var request device.LoadCertificates
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCertificates(ctx, dev, request)
	case "GetClientCertificateMode":
		return device_rpc.Call_GetClientCertificateMode(ctx, dev, device.GetClientCertificateMode{})
	case "SetClientCertificateMode":
		var request device.SetClientCertificateMode
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetClientCertificateMode(ctx, dev, request)
	case "GetRelayOutputs":
		return device_rpc.Call_GetRelayOutputs(ctx, dev, device.GetRelayOutputs{})
	case "SetRelayOutputSettings":
		var request device.SetRelayOutputSettings
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRelayOutputSettings(ctx, dev, request)
	case "SetRelayOutputState":
		var request device.SetRelayOutputState
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetRelayOutputState(ctx, dev, request)
	case "SendAuxiliaryCommand":
		var request device.SendAuxiliaryCommand
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SendAuxiliaryCommand(ctx, dev, request)
	case "GetCACertificates":
		return device_rpc.Call_GetCACertificates(ctx, dev, device.GetCACertificates{})
	case "LoadCertificateWithPrivateKey":
		var request device.LoadCertificateWithPrivateKey
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCertificateWithPrivateKey(ctx, dev, request)
	case "GetCertificateInformation":
		var request device.GetCertificateInformation
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetCertificateInformation(ctx, dev, request)
	case "LoadCACertificates":
		var request device.LoadCACertificates
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_LoadCACertificates(ctx, dev, request)
	case "CreateDot1XConfiguration":
		var request device.CreateDot1XConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateDot1XConfiguration(ctx, dev, request)
	case "SetDot1XConfiguration":
		var request device.SetDot1XConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetDot1XConfiguration(ctx, dev, request)
	case "GetDot1XConfiguration":
		var request device.GetDot1XConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetDot1XConfiguration(ctx, dev, request)
	case "GetDot1XConfigurations":
		return device_rpc.Call_GetDot1XConfigurations(ctx, dev, device.GetDot1XConfigurations{})
	case "DeleteDot1XConfiguration":
		var request device.DeleteDot1XConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteDot1XConfiguration(ctx, dev, request)
	case "GetDot11Capabilities":
		return device_rpc.Call_GetDot11Capabilities(ctx, dev, device.GetDot11Capabilities{})
	case "GetDot11Status":
		var request device.GetDot11Status
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetDot11Status(ctx, dev, request)
	case "ScanAvailableDot11Networks":
		var request device.ScanAvailableDot11Networks
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_ScanAvailableDot11Networks(ctx, dev, request)
	case "GetSystemUris":
		return device_rpc.Call_GetSystemUris(ctx, dev, device.GetSystemUris{})
	case "StartFirmwareUpgrade":
		return device_rpc.Call_StartFirmwareUpgrade(ctx, dev, device.StartFirmwareUpgrade{})
	case "StartSystemRestore":
		return device_rpc.Call_StartSystemRestore(ctx, dev, device.StartSystemRestore{})
	case "GetStorageConfigurations":
		return device_rpc.Call_GetStorageConfigurations(ctx, dev, device.GetStorageConfigurations{})
	case "CreateStorageConfiguration":
		var request device.CreateStorageConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_CreateStorageConfiguration(ctx, dev, request)
	case "GetStorageConfiguration":
		var request device.GetStorageConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_GetStorageConfiguration(ctx, dev, request)
	case "SetStorageConfiguration":
		var request device.SetStorageConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetStorageConfiguration(ctx, dev, request)
	case "DeleteStorageConfiguration":
		var request device.DeleteStorageConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteStorageConfiguration(ctx, dev, request)
	case "GetGeoLocation":
		return device_rpc.Call_GetGeoLocation(ctx, dev, device.GetGeoLocation{})
	case "SetGeoLocation":
		var request device.SetGeoLocation
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_SetGeoLocation(ctx, dev, request)
	case "DeleteGeoLocation":
		var request device.DeleteGeoLocation
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return device_rpc.Call_DeleteGeoLocation(ctx, dev, request)
	default:
		return nil, errors.New("unknown device method: " + methodName)
	}
}

// callPTZMethod handles all PTZ service methods
func callPTZMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
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
		var request ptz.GetPresetTours
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return ptz_rpc.Call_GetPresetTours(ctx, dev, request)
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
	default:
		return nil, errors.New("unknown PTZ method: " + methodName)
	}
}

// callMediaMethod handles all Media service methods
func callMediaMethod(methodName string, dev *onvif.Device, data []byte) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodName {
	case "GetServiceCapabilities":
		return media_rpc.Call_GetServiceCapabilities(ctx, dev, media.GetServiceCapabilities{})
	case "GetVideoSources":
		return media_rpc.Call_GetVideoSources(ctx, dev, media.GetVideoSources{})
	case "GetAudioSources":
		return media_rpc.Call_GetAudioSources(ctx, dev, media.GetAudioSources{})
	case "GetAudioOutputs":
		return media_rpc.Call_GetAudioOutputs(ctx, dev, media.GetAudioOutputs{})
	case "CreateProfile":
		var request media.CreateProfile
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_CreateProfile(ctx, dev, request)
	case "GetProfile":
		var request media.GetProfile
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetProfile(ctx, dev, request)
	case "GetProfiles":
		return media_rpc.Call_GetProfiles(ctx, dev, media.GetProfiles{})
	case "AddVideoEncoderConfiguration":
		var request media.AddVideoEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoEncoderConfiguration(ctx, dev, request)
	case "RemoveVideoEncoderConfiguration":
		var request media.RemoveVideoEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoEncoderConfiguration(ctx, dev, request)
	case "AddVideoSourceConfiguration":
		var request media.AddVideoSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoSourceConfiguration(ctx, dev, request)
	case "RemoveVideoSourceConfiguration":
		var request media.RemoveVideoSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoSourceConfiguration(ctx, dev, request)
	case "AddAudioEncoderConfiguration":
		var request media.AddAudioEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioEncoderConfiguration(ctx, dev, request)
	case "RemoveAudioEncoderConfiguration":
		var request media.RemoveAudioEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioEncoderConfiguration(ctx, dev, request)
	case "AddAudioSourceConfiguration":
		var request media.AddAudioSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioSourceConfiguration(ctx, dev, request)
	case "RemoveAudioSourceConfiguration":
		var request media.RemoveAudioSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioSourceConfiguration(ctx, dev, request)
	case "AddPTZConfiguration":
		var request media.AddPTZConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddPTZConfiguration(ctx, dev, request)
	case "RemovePTZConfiguration":
		var request media.RemovePTZConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemovePTZConfiguration(ctx, dev, request)
	case "AddVideoAnalyticsConfiguration":
		var request media.AddVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddVideoAnalyticsConfiguration(ctx, dev, request)
	case "RemoveVideoAnalyticsConfiguration":
		var request media.RemoveVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveVideoAnalyticsConfiguration(ctx, dev, request)
	case "AddMetadataConfiguration":
		var request media.AddMetadataConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddMetadataConfiguration(ctx, dev, request)
	case "RemoveMetadataConfiguration":
		var request media.RemoveMetadataConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveMetadataConfiguration(ctx, dev, request)
	case "AddAudioOutputConfiguration":
		var request media.AddAudioOutputConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioOutputConfiguration(ctx, dev, request)
	case "RemoveAudioOutputConfiguration":
		var request media.RemoveAudioOutputConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioOutputConfiguration(ctx, dev, request)
	case "AddAudioDecoderConfiguration":
		var request media.AddAudioDecoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_AddAudioDecoderConfiguration(ctx, dev, request)
	case "RemoveAudioDecoderConfiguration":
		var request media.RemoveAudioDecoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_RemoveAudioDecoderConfiguration(ctx, dev, request)
	case "DeleteProfile":
		var request media.DeleteProfile
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_DeleteProfile(ctx, dev, request)
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
		var request media.GetVideoSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoSourceConfiguration(ctx, dev, request)
	case "GetVideoEncoderConfiguration":
		var request media.GetVideoEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoEncoderConfiguration(ctx, dev, request)
	case "GetAudioSourceConfiguration":
		var request media.GetAudioSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioSourceConfiguration(ctx, dev, request)
	case "GetAudioEncoderConfiguration":
		var request media.GetAudioEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioEncoderConfiguration(ctx, dev, request)
	case "GetVideoAnalyticsConfiguration":
		var request media.GetVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoAnalyticsConfiguration(ctx, dev, request)
	case "GetMetadataConfiguration":
		var request media.GetMetadataConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetMetadataConfiguration(ctx, dev, request)
	case "GetAudioOutputConfiguration":
		var request media.GetAudioOutputConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioOutputConfiguration(ctx, dev, request)
	case "GetAudioDecoderConfiguration":
		var request media.GetAudioDecoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioDecoderConfiguration(ctx, dev, request)
	case "GetCompatibleVideoEncoderConfigurations":
		var request media.GetCompatibleVideoEncoderConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleVideoEncoderConfigurations(ctx, dev, request)
	case "GetCompatibleVideoSourceConfigurations":
		var request media.GetCompatibleVideoSourceConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleVideoSourceConfigurations(ctx, dev, request)
	case "GetCompatibleAudioEncoderConfigurations":
		var request media.GetCompatibleAudioEncoderConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleAudioEncoderConfigurations(ctx, dev, request)
	case "GetCompatibleAudioSourceConfigurations":
		var request media.GetCompatibleAudioSourceConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleAudioSourceConfigurations(ctx, dev, request)
	case "GetCompatibleVideoAnalyticsConfigurations":
		var request media.GetCompatibleVideoAnalyticsConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleVideoAnalyticsConfigurations(ctx, dev, request)
	case "GetCompatibleMetadataConfigurations":
		var request media.GetCompatibleMetadataConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleMetadataConfigurations(ctx, dev, request)
	case "GetCompatibleAudioOutputConfigurations":
		var request media.GetCompatibleAudioOutputConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleAudioOutputConfigurations(ctx, dev, request)
	case "GetCompatibleAudioDecoderConfigurations":
		var request media.GetCompatibleAudioDecoderConfigurations
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetCompatibleAudioDecoderConfigurations(ctx, dev, request)
	case "SetVideoSourceConfiguration":
		var request media.SetVideoSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoSourceConfiguration(ctx, dev, request)
	case "SetVideoEncoderConfiguration":
		var request media.SetVideoEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoEncoderConfiguration(ctx, dev, request)
	case "SetAudioSourceConfiguration":
		var request media.SetAudioSourceConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioSourceConfiguration(ctx, dev, request)
	case "SetAudioEncoderConfiguration":
		var request media.SetAudioEncoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioEncoderConfiguration(ctx, dev, request)
	case "SetVideoAnalyticsConfiguration":
		var request media.SetVideoAnalyticsConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoAnalyticsConfiguration(ctx, dev, request)
	case "SetMetadataConfiguration":
		var request media.SetMetadataConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetMetadataConfiguration(ctx, dev, request)
	case "SetAudioOutputConfiguration":
		var request media.SetAudioOutputConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioOutputConfiguration(ctx, dev, request)
	case "SetAudioDecoderConfiguration":
		var request media.SetAudioDecoderConfiguration
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetAudioDecoderConfiguration(ctx, dev, request)
	case "GetVideoSourceConfigurationOptions":
		var request media.GetVideoSourceConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoSourceConfigurationOptions(ctx, dev, request)
	case "GetVideoEncoderConfigurationOptions":
		var request media.GetVideoEncoderConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoEncoderConfigurationOptions(ctx, dev, request)
	case "GetAudioSourceConfigurationOptions":
		var request media.GetAudioSourceConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioSourceConfigurationOptions(ctx, dev, request)
	case "GetAudioEncoderConfigurationOptions":
		var request media.GetAudioEncoderConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioEncoderConfigurationOptions(ctx, dev, request)
	case "GetMetadataConfigurationOptions":
		var request media.GetMetadataConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetMetadataConfigurationOptions(ctx, dev, request)
	case "GetAudioOutputConfigurationOptions":
		var request media.GetAudioOutputConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioOutputConfigurationOptions(ctx, dev, request)
	case "GetAudioDecoderConfigurationOptions":
		var request media.GetAudioDecoderConfigurationOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetAudioDecoderConfigurationOptions(ctx, dev, request)
	case "GetGuaranteedNumberOfVideoEncoderInstances":
		var request media.GetGuaranteedNumberOfVideoEncoderInstances
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetGuaranteedNumberOfVideoEncoderInstances(ctx, dev, request)
	case "GetStreamUri":
		var request media.GetStreamUri
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetStreamUri(ctx, dev, request)
	case "StartMulticastStreaming":
		var request media.StartMulticastStreaming
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_StartMulticastStreaming(ctx, dev, request)
	case "StopMulticastStreaming":
		var request media.StopMulticastStreaming
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_StopMulticastStreaming(ctx, dev, request)
	case "SetSynchronizationPoint":
		var request media.SetSynchronizationPoint
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetSynchronizationPoint(ctx, dev, request)
	case "GetSnapshotUri":
		var request media.GetSnapshotUri
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetSnapshotUri(ctx, dev, request)
	case "GetVideoSourceModes":
		var request media.GetVideoSourceModes
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetVideoSourceModes(ctx, dev, request)
	case "SetVideoSourceMode":
		var request media.SetVideoSourceMode
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetVideoSourceMode(ctx, dev, request)
	case "GetOSDs":
		var request media.GetOSDs
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSDs(ctx, dev, request)
	case "GetOSD":
		var request media.GetOSD
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSD(ctx, dev, request)
	case "GetOSDOptions":
		var request media.GetOSDOptions
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_GetOSDOptions(ctx, dev, request)
	case "SetOSD":
		var request media.SetOSD
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_SetOSD(ctx, dev, request)
	case "CreateOSD":
		var request media.CreateOSD
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_CreateOSD(ctx, dev, request)
	case "DeleteOSD":
		var request media.DeleteOSD
		if err := json.Unmarshal(data, &request); err != nil {
			return nil, err
		}
		return media_rpc.Call_DeleteOSD(ctx, dev, request)
	default:
		return nil, errors.New("unknown Media method: " + methodName)
	}
}
