package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/rs/zerolog"

	"github.com/beevik/etree"
	"github.com/gin-gonic/gin"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/event"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	device_rpc "github.com/use-go/onvif/sdk/device"
	wsdiscovery "github.com/use-go/onvif/ws-discovery"
)

func getPTZStructByName(name string) (interface{}, error) {
	switch name {
	case "GetServiceCapabilities":
		return &ptz.GetServiceCapabilities{}, nil
	case "GetNodes":
		return &ptz.GetNodes{}, nil
	case "GetNode":
		return &ptz.GetNode{}, nil
	case "GetConfiguration":
		return &ptz.GetConfiguration{}, nil
	case "GetConfigurations":
		return &ptz.GetConfigurations{}, nil
	case "SetConfiguration":
		return &ptz.SetConfiguration{}, nil
	case "GetConfigurationOptions":
		return &ptz.GetConfigurationOptions{}, nil
	case "SendAuxiliaryCommand":
		return &ptz.SendAuxiliaryCommand{}, nil
	case "GetPresets":
		return &ptz.GetPresets{}, nil
	case "SetPreset":
		return &ptz.SetPreset{}, nil
	case "RemovePreset":
		return &ptz.RemovePreset{}, nil
	case "GotoPreset":
		return &ptz.GotoPreset{}, nil
	case "GotoHomePosition":
		return &ptz.GotoHomePosition{}, nil
	case "SetHomePosition":
		return &ptz.SetHomePosition{}, nil
	case "ContinuousMove":
		return &ptz.ContinuousMove{}, nil
	case "RelativeMove":
		return &ptz.RelativeMove{}, nil
	case "GetStatus":
		return &ptz.GetStatus{}, nil
	case "AbsoluteMove":
		return &ptz.AbsoluteMove{}, nil
	case "GeoMove":
		return &ptz.GeoMove{}, nil
	case "Stop":
		return &ptz.Stop{}, nil
	case "GetPresetTours":
		return &ptz.GetPresetTours{}, nil
	case "GetPresetTour":
		return &ptz.GetPresetTour{}, nil
	case "GetPresetTourOptions":
		return &ptz.GetPresetTourOptions{}, nil
	case "CreatePresetTour":
		return &ptz.CreatePresetTour{}, nil
	case "ModifyPresetTour":
		return &ptz.ModifyPresetTour{}, nil
	case "OperatePresetTour":
		return &ptz.OperatePresetTour{}, nil
	case "RemovePresetTour":
		return &ptz.RemovePresetTour{}, nil
	case "GetCompatibleConfigurations":
		return &ptz.GetCompatibleConfigurations{}, nil
	default:
		return nil, errors.New("there is no such method in the PTZ service")
	}
}

func getDeviceStructByName(name string, dev *onvif.Device, acceptedData []byte) (interface{}, error) {
	// Use a context with timeout for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch name {
	// case "GetServices":
	// 	return &device.GetServices{}, nil
	// case "GetServiceCapabilities":
	// 	return &device.GetServiceCapabilities{}, nil
	case "GetDeviceInformation":
		return device_rpc.Call_GetDeviceInformation(ctx, dev, device.GetDeviceInformation{})
	// case "SetSystemDateAndTime":
	// 	return &device.SetSystemDateAndTime{}, nil
	case "GetSystemDateAndTime":
		return device_rpc.Call_GetSystemDateAndTime(ctx, dev, device.GetSystemDateAndTime{})
	// case "SetSystemFactoryDefault":
	// 	return &device.SetSystemFactoryDefault{}, nil
	// case "UpgradeSystemFirmware":
	// 	return &device.UpgradeSystemFirmware{}, nil
	// case "SystemReboot":
	// 	return &device.SystemReboot{}, nil
	// case "RestoreSystem":
	// 	return &device.RestoreSystem{}, nil
	case "GetSystemBackup":
		return device_rpc.Call_GetSystemBackup(ctx, dev, device.GetSystemBackup{})
	case "GetSystemLog":
		return device_rpc.Call_GetSystemLog(ctx, dev, device.GetSystemLog{})
	case "GetSystemSupportInformation":
		return device_rpc.Call_GetSystemSupportInformation(ctx, dev, device.GetSystemSupportInformation{})
	case "GetScopes":
		return device_rpc.Call_GetScopes(ctx, dev, device.GetScopes{})
	// case "SetScopes":
	// 	return &device.SetScopes{}, nil
	// case "AddScopes":
	// 	return &device.AddScopes{}, nil
	// case "RemoveScopes":
	// 	return &device.RemoveScopes{}, nil
	case "GetDiscoveryMode":
		return device_rpc.Call_GetDiscoveryMode(ctx, dev, device.GetDiscoveryMode{})
	// case "SetDiscoveryMode":
	// 	return &device.SetDiscoveryMode{}, nil
	case "GetRemoteDiscoveryMode":
		return device_rpc.Call_GetRemoteDiscoveryMode(ctx, dev, device.GetRemoteDiscoveryMode{})
	// case "SetRemoteDiscoveryMode":
	// 	return &device.SetRemoteDiscoveryMode{}, nil
	case "GetDPAddresses":
		return device_rpc.Call_GetDPAddresses(ctx, dev, device.GetDPAddresses{})
	// case "SetDPAddresses":
	// 	return &device.SetDPAddresses{}, nil
	case "GetEndpointReference":
		return device_rpc.Call_GetEndpointReference(ctx, dev, device.GetEndpointReference{})
	case "GetRemoteUser":
		return device_rpc.Call_GetRemoteUser(ctx, dev, device.GetRemoteUser{})
	// case "SetRemoteUser":
	// 	return &device.SetRemoteUser{}, nil
	case "GetUsers":
		return device_rpc.Call_GetUsers(ctx, dev, device.GetUsers{})
	// case "CreateUsers":
	// 	return &device.CreateUsers{}, nil
	// case "DeleteUsers":
	// 	return &device.DeleteUsers{}, nil
	// case "SetUser":
	// 	return &device.SetUser{}, nil
	case "GetWsdlUrl":
		return device_rpc.Call_GetWsdlUrl(ctx, dev, device.GetWsdlUrl{})
	case "GetCapabilities":
		var request device.GetCapabilities

		// Unmarshal JSON into struct
		if err := json.Unmarshal(acceptedData, &request); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			return "", err
		}

		return device_rpc.Call_GetCapabilities(ctx, dev, request)

	case "GetHostname":
		return &device.GetHostname{}, nil
	// case "SetHostname":
	// 	return &device.SetHostname{}, nil
	// case "SetHostnameFromDHCP":
	// 	return &device.SetHostnameFromDHCP{}, nil
	case "GetDNS":
		return device_rpc.Call_GetDNS(ctx, dev, device.GetDNS{})
	// case "SetDNS":
	// 	return &device.SetDNS{}, nil
	case "GetNTP":
		return device_rpc.Call_GetNTP(ctx, dev, device.GetNTP{})
	// case "SetNTP":
	// 	return &device.SetNTP{}, nil
	case "GetDynamicDNS":
		return device_rpc.Call_GetDynamicDNS(ctx, dev, device.GetDynamicDNS{})
	// case "SetDynamicDNS":
	// 	return &device.SetDynamicDNS{}, nil
	case "GetNetworkInterfaces":
		return device_rpc.Call_GetNetworkInterfaces(ctx, dev, device.GetNetworkInterfaces{})
	// case "SetNetworkInterfaces":
	// 	return &device.SetNetworkInterfaces{}, nil
	case "GetNetworkProtocols":
		return device_rpc.Call_GetNetworkProtocols(ctx, dev, device.GetNetworkProtocols{})
	// case "SetNetworkProtocols":
	// 	return &device.SetNetworkProtocols{}, nil
	case "GetNetworkDefaultGateway":
		return device_rpc.Call_GetNetworkDefaultGateway(ctx, dev, device.GetNetworkDefaultGateway{})
	// case "SetNetworkDefaultGateway":
	// 	return &device.SetNetworkDefaultGateway{}, nil
	case "GetZeroConfiguration":
		return device_rpc.Call_GetZeroConfiguration(ctx, dev, device.GetZeroConfiguration{})
	// case "SetZeroConfiguration":
	// 	return &device.SetZeroConfiguration{}, nil
	case "GetIPAddressFilter":
		return device_rpc.Call_GetIPAddressFilter(ctx, dev, device.GetIPAddressFilter{})
	// case "SetIPAddressFilter":
	// 	return &device.SetIPAddressFilter{}, nil
	// case "AddIPAddressFilter":
	// 	return &device.AddIPAddressFilter{}, nil
	// case "RemoveIPAddressFilter":
	// 	return &device.RemoveIPAddressFilter{}, nil
	case "GetAccessPolicy":
		return device_rpc.Call_GetAccessPolicy(ctx, dev, device.GetAccessPolicy{})
	// case "SetAccessPolicy":
	// 	return &device.SetAccessPolicy{}, nil
	// case "CreateCertificate":
	// 	return &device.CreateCertificate{}, nil
	case "GetCertificates":
		return device_rpc.Call_GetCertificates(ctx, dev, device.GetCertificates{})
	case "GetCertificatesStatus":
		return device_rpc.Call_GetCertificatesStatus(ctx, dev, device.GetCertificatesStatus{})
	// case "SetCertificatesStatus":
	// 	return &device.SetCertificatesStatus{}, nil
	// case "DeleteCertificates":
	// 	return &device.DeleteCertificates{}, nil
	case "GetPkcs10Request":
		return device_rpc.Call_GetPkcs10Request(ctx, dev, device.GetPkcs10Request{})
	// case "LoadCertificates":
	// 	return &device.LoadCertificates{}, nil
	case "GetClientCertificateMode":
		return device_rpc.Call_GetClientCertificateMode(ctx, dev, device.GetClientCertificateMode{})
	// case "SetClientCertificateMode":
	// 	return &device.SetClientCertificateMode{}, nil
	case "GetRelayOutputs":
		return device_rpc.Call_GetRelayOutputs(ctx, dev, device.GetRelayOutputs{})
	// case "SetRelayOutputSettings":
	// 	return &device.SetRelayOutputSettings{}, nil
	// case "SetRelayOutputState":
	// 	return &device.SetRelayOutputState{}, nil
	// case "SendAuxiliaryCommand":
	// 	return &device.SendAuxiliaryCommand{}, nil
	case "GetCACertificates":
		return device_rpc.Call_GetCACertificates(ctx, dev, device.GetCACertificates{})
	// case "LoadCertificateWithPrivateKey":
	// 	return &device.LoadCertificateWithPrivateKey{}, nil
	case "GetCertificateInformation":
		return device_rpc.Call_GetCertificateInformation(ctx, dev, device.GetCertificateInformation{})
	// case "LoadCACertificates":
	// 	return &device.LoadCACertificates{}, nil
	// case "CreateDot1XConfiguration":
	// 	return &device.CreateDot1XConfiguration{}, nil
	// case "SetDot1XConfiguration":
	// 	return &device.SetDot1XConfiguration{}, nil
	case "GetDot1XConfiguration":
		return device_rpc.Call_GetDot1XConfiguration(ctx, dev, device.GetDot1XConfiguration{})
	case "GetDot1XConfigurations":
		return device_rpc.Call_GetDot1XConfigurations(ctx, dev, device.GetDot1XConfigurations{})
	// case "DeleteDot1XConfiguration":
	// 	return &device.DeleteDot1XConfiguration{}, nil
	case "GetDot11Capabilities":
		return device_rpc.Call_GetDot11Capabilities(ctx, dev, device.GetDot11Capabilities{})
	case "GetDot11Status":
		return device_rpc.Call_GetDot11Status(ctx, dev, device.GetDot11Status{})
	// case "ScanAvailableDot11Networks":
	// 	return &device.ScanAvailableDot11Networks{}, nil
	case "GetSystemUris":
		return device_rpc.Call_GetSystemUris(ctx, dev, device.GetSystemUris{})
	// case "StartFirmwareUpgrade":
	// 	return &device.StartFirmwareUpgrade{}, nil
	// case "StartSystemRestore":
	// 	return &device.StartSystemRestore{}, nil
	case "GetStorageConfigurations":
		return device_rpc.Call_GetStorageConfigurations(ctx, dev, device.GetStorageConfigurations{})
	// case "CreateStorageConfiguration":
	// 	return &device.CreateStorageConfiguration{}, nil
	case "GetStorageConfiguration":
		return &device.GetStorageConfiguration{}, nil
	// case "SetStorageConfiguration":
	// 	return &device.SetStorageConfiguration{}, nil
	// case "DeleteStorageConfiguration":
	// 	return &device.DeleteStorageConfiguration{}, nil
	case "GetGeoLocation":
		return device_rpc.Call_GetGeoLocation(ctx, dev, device.GetGeoLocation{})
	// case "SetGeoLocation":
	// 	return &device.SetGeoLocation{}, nil
	// case "DeleteGeoLocation":
	// 	return &device.DeleteGeoLocation{}, nil
	default:
		return "", errors.New("there is no such method in the Device service")
	}
}

func getMediaStructByName(name string) (interface{}, error) {
	switch name {
	case "GetServiceCapabilities":
		return &media.GetServiceCapabilities{}, nil
	case "GetVideoSources":
		return &media.GetVideoSources{}, nil
	case "GetAudioSources":
		return &media.GetAudioSources{}, nil
	case "GetAudioOutputs":
		return &media.GetAudioOutputs{}, nil
	case "CreateProfile":
		return &media.CreateProfile{}, nil
	case "GetProfile":
		return &media.GetProfile{}, nil
	case "GetProfiles":
		return &media.GetProfiles{}, nil
	case "AddVideoEncoderConfiguration":
		return &media.AddVideoEncoderConfiguration{}, nil
	case "RemoveVideoEncoderConfiguration":
		return &media.RemoveVideoEncoderConfiguration{}, nil
	case "AddVideoSourceConfiguration":
		return &media.AddVideoSourceConfiguration{}, nil
	case "RemoveVideoSourceConfiguration":
		return &media.RemoveVideoSourceConfiguration{}, nil
	case "AddAudioEncoderConfiguration":
		return &media.AddAudioEncoderConfiguration{}, nil
	case "RemoveAudioEncoderConfiguration":
		return &media.RemoveAudioEncoderConfiguration{}, nil
	case "AddAudioSourceConfiguration":
		return &media.AddAudioSourceConfiguration{}, nil
	case "RemoveAudioSourceConfiguration":
		return &media.RemoveAudioSourceConfiguration{}, nil
	case "AddPTZConfiguration":
		return &media.AddPTZConfiguration{}, nil
	case "RemovePTZConfiguration":
		return &media.RemovePTZConfiguration{}, nil
	case "AddVideoAnalyticsConfiguration":
		return &media.AddVideoAnalyticsConfiguration{}, nil
	case "RemoveVideoAnalyticsConfiguration":
		return &media.RemoveVideoAnalyticsConfiguration{}, nil
	case "AddMetadataConfiguration":
		return &media.AddMetadataConfiguration{}, nil
	case "RemoveMetadataConfiguration":
		return &media.RemoveMetadataConfiguration{}, nil
	case "AddAudioOutputConfiguration":
		return &media.AddAudioOutputConfiguration{}, nil
	case "RemoveAudioOutputConfiguration":
		return &media.RemoveAudioOutputConfiguration{}, nil
	case "AddAudioDecoderConfiguration":
		return &media.AddAudioDecoderConfiguration{}, nil
	case "RemoveAudioDecoderConfiguration":
		return &media.RemoveAudioDecoderConfiguration{}, nil
	case "DeleteProfile":
		return &media.DeleteProfile{}, nil
	case "GetVideoSourceConfigurations":
		return &media.GetVideoSourceConfigurations{}, nil
	case "GetVideoEncoderConfigurations":
		return &media.GetVideoEncoderConfigurations{}, nil
	case "GetAudioSourceConfigurations":
		return &media.GetAudioSourceConfigurations{}, nil
	case "GetAudioEncoderConfigurations":
		return &media.GetAudioEncoderConfigurations{}, nil
	case "GetVideoAnalyticsConfigurations":
		return &media.GetVideoAnalyticsConfigurations{}, nil
	case "GetMetadataConfigurations":
		return &media.GetMetadataConfigurations{}, nil
	case "GetAudioOutputConfigurations":
		return &media.GetAudioOutputConfigurations{}, nil
	case "GetAudioDecoderConfigurations":
		return &media.GetAudioDecoderConfigurations{}, nil
	case "GetVideoSourceConfiguration":
		return &media.GetVideoSourceConfiguration{}, nil
	case "GetVideoEncoderConfiguration":
		return &media.GetVideoEncoderConfiguration{}, nil
	case "GetAudioSourceConfiguration":
		return &media.GetAudioSourceConfiguration{}, nil
	case "GetAudioEncoderConfiguration":
		return &media.GetAudioEncoderConfiguration{}, nil
	case "GetVideoAnalyticsConfiguration":
		return &media.GetVideoAnalyticsConfiguration{}, nil
	case "GetMetadataConfiguration":
		return &media.GetMetadataConfiguration{}, nil
	case "GetAudioOutputConfiguration":
		return &media.GetAudioOutputConfiguration{}, nil
	case "GetAudioDecoderConfiguration":
		return &media.GetAudioDecoderConfiguration{}, nil
	case "GetCompatibleVideoEncoderConfigurations":
		return &media.GetCompatibleVideoEncoderConfigurations{}, nil
	case "GetCompatibleVideoSourceConfigurations":
		return &media.GetCompatibleVideoSourceConfigurations{}, nil
	case "GetCompatibleAudioEncoderConfigurations":
		return &media.GetCompatibleAudioEncoderConfigurations{}, nil
	case "GetCompatibleAudioSourceConfigurations":
		return &media.GetCompatibleAudioSourceConfigurations{}, nil
	case "GetCompatibleVideoAnalyticsConfigurations":
		return &media.GetCompatibleVideoAnalyticsConfigurations{}, nil
	case "GetCompatibleMetadataConfigurations":
		return &media.GetCompatibleMetadataConfigurations{}, nil
	case "GetCompatibleAudioOutputConfigurations":
		return &media.GetCompatibleAudioOutputConfigurations{}, nil
	case "GetCompatibleAudioDecoderConfigurations":
		return &media.GetCompatibleAudioDecoderConfigurations{}, nil
	case "SetVideoSourceConfiguration":
		return &media.SetVideoSourceConfiguration{}, nil
	case "SetVideoEncoderConfiguration":
		return &media.SetVideoEncoderConfiguration{}, nil
	case "SetAudioSourceConfiguration":
		return &media.SetAudioSourceConfiguration{}, nil
	case "SetAudioEncoderConfiguration":
		return &media.SetAudioEncoderConfiguration{}, nil
	case "SetVideoAnalyticsConfiguration":
		return &media.SetVideoAnalyticsConfiguration{}, nil
	case "SetMetadataConfiguration":
		return &media.SetMetadataConfiguration{}, nil
	case "SetAudioOutputConfiguration":
		return &media.SetAudioOutputConfiguration{}, nil
	case "SetAudioDecoderConfiguration":
		return &media.SetAudioDecoderConfiguration{}, nil
	case "GetVideoSourceConfigurationOptions":
		return &media.GetVideoSourceConfigurationOptions{}, nil
	case "GetVideoEncoderConfigurationOptions":
		return &media.GetVideoEncoderConfigurationOptions{}, nil
	case "GetAudioSourceConfigurationOptions":
		return &media.GetAudioSourceConfigurationOptions{}, nil
	case "GetAudioEncoderConfigurationOptions":
		return &media.GetAudioEncoderConfigurationOptions{}, nil
	case "GetMetadataConfigurationOptions":
		return &media.GetMetadataConfigurationOptions{}, nil
	case "GetAudioOutputConfigurationOptions":
		return &media.GetAudioOutputConfigurationOptions{}, nil
	case "GetAudioDecoderConfigurationOptions":
		return &media.GetAudioDecoderConfigurationOptions{}, nil
	case "GetGuaranteedNumberOfVideoEncoderInstances":
		return &media.GetGuaranteedNumberOfVideoEncoderInstances{}, nil
	case "GetStreamUri":
		return &media.GetStreamUri{}, nil
	case "StartMulticastStreaming":
		return &media.StartMulticastStreaming{}, nil
	case "StopMulticastStreaming":
		return &media.StopMulticastStreaming{}, nil
	case "SetSynchronizationPoint":
		return &media.SetSynchronizationPoint{}, nil
	case "GetSnapshotUri":
		return &media.GetSnapshotUri{}, nil
	case "GetVideoSourceModes":
		return &media.GetVideoSourceModes{}, nil
	case "SetVideoSourceMode":
		return &media.SetVideoSourceMode{}, nil
	case "GetOSDs":
		return &media.GetOSDs{}, nil
	case "GetOSD":
		return &media.GetOSD{}, nil
	case "GetOSDOptions":
		return &media.GetOSDOptions{}, nil
	case "SetOSD":
		return &media.SetOSD{}, nil
	case "CreateOSD":
		return &media.CreateOSD{}, nil
	case "DeleteOSD":
		return &media.DeleteOSD{}, nil
	default:
		return nil, errors.New("there is no such method in the Media service")
	}

}

func getEventStructByName(name string) (interface{}, error) {
	switch name {
	case "GetEventProperties":
		return &event.GetEventProperties{}, nil
	case "CreatePullPointSubscription":
		return &event.CreatePullPointSubscription{}, nil
	case "Renew":
		return &event.Renew{}, nil
	case "Unsubscribe":
		return &event.Unsubscribe{}, nil
	case "GetServiceCapabilities":
		return &event.GetServiceCapabilities{}, nil
	default:
		return nil, errors.New("there is no such method in the Event service")
	}
}

var (
	// LoggerContext is the builder of a zerolog.Logger that is exposed to the application so that
	// options at the CLI might alter the formatting and the output of the logs.
	LoggerContext = zerolog.
			New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp()

	// Logger is a zerolog logger, that can be safely used from any part of the application.
	// It gathers the format and the output.
	Logger = LoggerContext.Logger()
)

func RunApi() {
	router := gin.Default()

	router.POST("/:service/:method", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.Header("Content-Type", "application/json")

		serviceName := c.Param("service")
		methodName := c.Param("method")
		username := c.GetHeader("username")
		pass := c.GetHeader("password")
		xaddr := c.GetHeader("xaddr")
		acceptedData, err := c.GetRawData()
		if err != nil {
			Logger.Debug().Err(err).Msg("Failed to get rawx data")
		}

		dev, err := onvif.NewDevice(onvif.DeviceParams{
			Xaddr:    xaddr,
			Username: username,
			Password: pass,
		})

		if err != nil {
			c.XML(http.StatusBadRequest, err.Error())
		}

		message, err := callNecessaryMethod(serviceName, methodName, []byte(acceptedData), dev)
		if err != nil {
			c.XML(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusOK, message)
		}
	})

	router.GET("/discovery", func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		interfaceName := context.GetHeader("interface")

		devices, err := wsdiscovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
		if err != nil {
			context.String(http.StatusInternalServerError, "error")
		} else {
			response := "["

			for _, j := range devices {
				doc := etree.NewDocument()

				if err := doc.ReadFromString(j); err != nil {
					context.XML(http.StatusBadRequest, err.Error())
				} else {

					endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
					scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")

					flag := false

					for _, xaddr := range endpoints {
						xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
						if strings.Contains(response, xaddr) {
							flag = true
							break
						}
						response += "{"
						response += `"url":"` + xaddr + `",`
					}
					if flag {
						break
					}
					for _, scope := range scopes {
						re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
						match := re.FindStringSubmatch(scope.Text())
						response += `"name":"` + path.Base(match[0]) + `"`
					}
					response += "},"
				}
			}
			response = strings.TrimRight(response, ",")
			response += "]"
			context.String(http.StatusOK, response)
		}
	})

	router.Run()
}

func callNecessaryMethod(serviceName, methodName string, acceptedData []byte, dev *onvif.Device) (string, error) {
	var response interface{}
	var err error

	fmt.Println(acceptedData)

	switch strings.ToLower(serviceName) {
	case "device":
		response, err = getDeviceStructByName(methodName, dev, acceptedData)
	// case "ptz":
	// 	methodStruct, err = getPTZStructByName(methodName)
	// case "media":
	// 	methodStruct, err = getMediaStructByName(methodName)
	// case "event":
	// 	methodStruct, err = getEventStructByName(methodName)
	default:
		return "", errors.New("there is no such service")
	}
	if err != nil { //done
		return "", errors.Annotate(err, "getStructByName")
	}

	// Marshal response to JSON for consistency
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	RunApi()
}
