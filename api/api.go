package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
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
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/networking"
	"github.com/use-go/onvif/ptz"

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

func getDeviceStructByName(name string) (interface{}, error) {
	switch name {
	case "GetServices":
		return &device.GetServices{}, nil
	case "GetServiceCapabilities":
		return &device.GetServiceCapabilities{}, nil
	case "GetDeviceInformation":
		return &device.GetDeviceInformation{}, nil
	case "SetSystemDateAndTime":
		return &device.SetSystemDateAndTime{}, nil
	case "GetSystemDateAndTime":
		return &device.GetSystemDateAndTime{}, nil
	case "SetSystemFactoryDefault":
		return &device.SetSystemFactoryDefault{}, nil
	case "UpgradeSystemFirmware":
		return &device.UpgradeSystemFirmware{}, nil
	case "SystemReboot":
		return &device.SystemReboot{}, nil
	case "RestoreSystem":
		return &device.RestoreSystem{}, nil
	case "GetSystemBackup":
		return &device.GetSystemBackup{}, nil
	case "GetSystemLog":
		return &device.GetSystemLog{}, nil
	case "GetSystemSupportInformation":
		return &device.GetSystemSupportInformation{}, nil
	case "GetScopes":
		return &device.GetScopes{}, nil
	case "SetScopes":
		return &device.SetScopes{}, nil
	case "AddScopes":
		return &device.AddScopes{}, nil
	case "RemoveScopes":
		return &device.RemoveScopes{}, nil
	case "GetDiscoveryMode":
		return &device.GetDiscoveryMode{}, nil
	case "SetDiscoveryMode":
		return &device.SetDiscoveryMode{}, nil
	case "GetRemoteDiscoveryMode":
		return &device.GetRemoteDiscoveryMode{}, nil
	case "SetRemoteDiscoveryMode":
		return &device.SetRemoteDiscoveryMode{}, nil
	case "GetDPAddresses":
		return &device.GetDPAddresses{}, nil
	case "SetDPAddresses":
		return &device.SetDPAddresses{}, nil
	case "GetEndpointReference":
		return &device.GetEndpointReference{}, nil
	case "GetRemoteUser":
		return &device.GetRemoteUser{}, nil
	case "SetRemoteUser":
		return &device.SetRemoteUser{}, nil
	case "GetUsers":
		return &device.GetUsers{}, nil
	case "CreateUsers":
		return &device.CreateUsers{}, nil
	case "DeleteUsers":
		return &device.DeleteUsers{}, nil
	case "SetUser":
		return &device.SetUser{}, nil
	case "GetWsdlUrl":
		return &device.GetWsdlUrl{}, nil
	case "GetCapabilities":
		return &device.GetCapabilities{}, nil
	case "GetHostname":
		return &device.GetHostname{}, nil
	case "SetHostname":
		return &device.SetHostname{}, nil
	case "SetHostnameFromDHCP":
		return &device.SetHostnameFromDHCP{}, nil
	case "GetDNS":
		return &device.GetDNS{}, nil
	case "SetDNS":
		return &device.SetDNS{}, nil
	case "GetNTP":
		return &device.GetNTP{}, nil
	case "SetNTP":
		return &device.SetNTP{}, nil
	case "GetDynamicDNS":
		return &device.GetDynamicDNS{}, nil
	case "SetDynamicDNS":
		return &device.SetDynamicDNS{}, nil
	case "GetNetworkInterfaces":
		return &device.GetNetworkInterfaces{}, nil
	case "SetNetworkInterfaces":
		return &device.SetNetworkInterfaces{}, nil
	case "GetNetworkProtocols":
		return &device.GetNetworkProtocols{}, nil
	case "SetNetworkProtocols":
		return &device.SetNetworkProtocols{}, nil
	case "GetNetworkDefaultGateway":
		return &device.GetNetworkDefaultGateway{}, nil
	case "SetNetworkDefaultGateway":
		return &device.SetNetworkDefaultGateway{}, nil
	case "GetZeroConfiguration":
		return &device.GetZeroConfiguration{}, nil
	case "SetZeroConfiguration":
		return &device.SetZeroConfiguration{}, nil
	case "GetIPAddressFilter":
		return &device.GetIPAddressFilter{}, nil
	case "SetIPAddressFilter":
		return &device.SetIPAddressFilter{}, nil
	case "AddIPAddressFilter":
		return &device.AddIPAddressFilter{}, nil
	case "RemoveIPAddressFilter":
		return &device.RemoveIPAddressFilter{}, nil
	case "GetAccessPolicy":
		return &device.GetAccessPolicy{}, nil
	case "SetAccessPolicy":
		return &device.SetAccessPolicy{}, nil
	case "CreateCertificate":
		return &device.CreateCertificate{}, nil
	case "GetCertificates":
		return &device.GetCertificates{}, nil
	case "GetCertificatesStatus":
		return &device.GetCertificatesStatus{}, nil
	case "SetCertificatesStatus":
		return &device.SetCertificatesStatus{}, nil
	case "DeleteCertificates":
		return &device.DeleteCertificates{}, nil
	case "GetPkcs10Request":
		return &device.GetPkcs10Request{}, nil
	case "LoadCertificates":
		return &device.LoadCertificates{}, nil
	case "GetClientCertificateMode":
		return &device.GetClientCertificateMode{}, nil
	case "SetClientCertificateMode":
		return &device.SetClientCertificateMode{}, nil
	case "GetRelayOutputs":
		return &device.GetRelayOutputs{}, nil
	case "SetRelayOutputSettings":
		return &device.SetRelayOutputSettings{}, nil
	case "SetRelayOutputState":
		return &device.SetRelayOutputState{}, nil
	case "SendAuxiliaryCommand":
		return &device.SendAuxiliaryCommand{}, nil
	case "GetCACertificates":
		return &device.GetCACertificates{}, nil
	case "LoadCertificateWithPrivateKey":
		return &device.LoadCertificateWithPrivateKey{}, nil
	case "GetCertificateInformation":
		return &device.GetCertificateInformation{}, nil
	case "LoadCACertificates":
		return &device.LoadCACertificates{}, nil
	case "CreateDot1XConfiguration":
		return &device.CreateDot1XConfiguration{}, nil
	case "SetDot1XConfiguration":
		return &device.SetDot1XConfiguration{}, nil
	case "GetDot1XConfiguration":
		return &device.GetDot1XConfiguration{}, nil
	case "GetDot1XConfigurations":
		return &device.GetDot1XConfigurations{}, nil
	case "DeleteDot1XConfiguration":
		return &device.DeleteDot1XConfiguration{}, nil
	case "GetDot11Capabilities":
		return &device.GetDot11Capabilities{}, nil
	case "GetDot11Status":
		return &device.GetDot11Status{}, nil
	case "ScanAvailableDot11Networks":
		return &device.ScanAvailableDot11Networks{}, nil
	case "GetSystemUris":
		return &device.GetSystemUris{}, nil
	case "StartFirmwareUpgrade":
		return &device.StartFirmwareUpgrade{}, nil
	case "StartSystemRestore":
		return &device.StartSystemRestore{}, nil
	case "GetStorageConfigurations":
		return &device.GetStorageConfigurations{}, nil
	case "CreateStorageConfiguration":
		return &device.CreateStorageConfiguration{}, nil
	case "GetStorageConfiguration":
		return &device.GetStorageConfiguration{}, nil
	case "SetStorageConfiguration":
		return &device.SetStorageConfiguration{}, nil
	case "DeleteStorageConfiguration":
		return &device.DeleteStorageConfiguration{}, nil
	case "GetGeoLocation":
		return &device.GetGeoLocation{}, nil
	case "SetGeoLocation":
		return &device.SetGeoLocation{}, nil
	case "DeleteGeoLocation":
		return &device.DeleteGeoLocation{}, nil
	default:
		return nil, errors.New("there is no such method in the Device service")
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

func getEventStructByName(name string) (interface {}, error) {
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
		c.Header("Content-Type", "application/xml")

		serviceName := c.Param("service")
		methodName := c.Param("method")
		username := c.GetHeader("username")
		pass := c.GetHeader("password")
		xaddr := c.GetHeader("xaddr")
		acceptedData, err := c.GetRawData()
		if err != nil {
			Logger.Debug().Err(err).Msg("Failed to get rawx data")
		}

		fmt.Println("serviceName", serviceName)
		fmt.Println("methodName", methodName)
		fmt.Println("username", username)
		fmt.Println("pass", pass)
		fmt.Println("xaddr", xaddr)
		fmt.Println("acceptedData", string(acceptedData))
		message, err := callNecessaryMethod(serviceName, methodName, string(acceptedData), username, pass, xaddr)
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

func callNecessaryMethod(serviceName, methodName, acceptedData, username, password, xaddr string) (string, error) {
	var methodStruct interface{}
	var err error

	switch strings.ToLower(serviceName) {
	case "device":
		methodStruct, err = getDeviceStructByName(methodName)
	case "ptz":
		methodStruct, err = getPTZStructByName(methodName)
	case "media":
		methodStruct, err = getMediaStructByName(methodName)
	case "event":
		methodStruct, err = getEventStructByName(methodName) 
	default:
		return "", errors.New("there is no such service")
	}
	if err != nil { //done
		return "", errors.Annotate(err, "getStructByName")
	}

	fmt.Println("methodStruct", methodStruct)
	fmt.Println("acceptedData", acceptedData)
	resp, err := xmlAnalize(methodStruct, &acceptedData)
	fmt.Println("resp", resp)
	if err != nil {
		return "", errors.Annotate(err, "xmlAnalize")
	}

	endpoint, err := getEndpoint(serviceName, xaddr)
	fmt.Println("endpoint", endpoint)
	if err != nil {
		return "", errors.Annotate(err, "getEndpoint")
	}

	soap := gosoap.NewEmptySOAP()
	soap.AddStringBodyContent(*resp)
	soap.AddRootNamespaces(onvif.Xlmns)
	soap.AddWSSecurity(username, password)

	fmt.Println("soap", soap)

	servResp, err := networking.SendSoap(new(http.Client), endpoint, soap.String())
	fmt.Println("servResp", servResp)
	if err != nil {
		return "", errors.Annotate(err, "SendSoap")
	}

	rsp, err := ioutil.ReadAll(servResp.Body)
	if err != nil {
		return "", errors.Annotate(err, "ReadAll")
	}

	servResp.Body.Close()
	

	return string(rsp), nil
}

func getEndpoint(service, xaddr string) (string, error) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: xaddr})
	fmt.Println("dev", dev)
	if err != nil {
		return "", errors.Annotate(err, "NewDevice")
	}
	pkg := strings.ToLower(service)
	fmt.Println("pkg", pkg)

	var endpoint string
	switch pkg {
	case "device":
		endpoint = dev.GetEndpoint("Device")
	case "event":
		endpoint = dev.GetEndpoint("Event")
	case "imaging":
		endpoint = dev.GetEndpoint("Imaging")
	case "media":
		endpoint = dev.GetEndpoint("Media")
	case "ptz":
		endpoint = dev.GetEndpoint("PTZ")
	}
	return endpoint, nil
}

func xmlAnalize(methodStruct interface{}, acceptedData *string) (*string, error) {
	test := make([]map[string]string, 0)      //tags
	testunMarshal := make([][]interface{}, 0) //data
	var mas []string                          //idnt

	soapHandling(methodStruct, &test)
	test = mapProcessing(test)

	doc := etree.NewDocument()
	if err := doc.ReadFromString(*acceptedData); err != nil {
		return nil, errors.Annotate(err, "readFromString")
	}
	etr := doc.FindElements("./*")
	xmlUnmarshal(etr, &testunMarshal, &mas)
	ident(&mas)

	document := etree.NewDocument()
	var el *etree.Element
	var idntIndex = 0

	for lstIndex := 0; lstIndex < len(testunMarshal); {
		lst := (testunMarshal)[lstIndex]
		elemName, attr, value, err := xmlMaker(&lst, &test, lstIndex)
		if err != nil {
			return nil, errors.Annotate(err, "xmlMarker")
		}

		if mas[lstIndex] == "Push" && lstIndex == 0 { //done
			el = document.CreateElement(elemName)
			el.SetText(value)
			if len(attr) != 0 {
				for key, value := range attr {
					el.CreateAttr(key, value)
				}
			}
		} else if mas[idntIndex] == "Push" {
			pushTmp := etree.NewElement(elemName)
			pushTmp.SetText(value)
			if len(attr) != 0 {
				for key, value := range attr {
					pushTmp.CreateAttr(key, value)
				}
			}
			el.AddChild(pushTmp)
			el = pushTmp
		} else if mas[idntIndex] == "PushPop" {
			popTmp := etree.NewElement(elemName)
			popTmp.SetText(value)
			if len(attr) != 0 {
				for key, value := range attr {
					popTmp.CreateAttr(key, value)
				}
			}
			if el == nil {
				document.AddChild(popTmp)
			} else {
				el.AddChild(popTmp)
			}
		} else if mas[idntIndex] == "Pop" {
			el = el.Parent()
			lstIndex -= 1
		}
		idntIndex += 1
		lstIndex += 1
	}

	resp, err := document.WriteToString()
	if err != nil {
		return nil, errors.Annotate(err, "writeToString")
	}

	return &resp, nil
}

func xmlMaker(lst *[]interface{}, tags *[]map[string]string, lstIndex int) (string, map[string]string, string, error) {
	var elemName, value string
	attr := make(map[string]string)
	for tgIndx, tg := range *tags {
		if tgIndx == lstIndex {
			for index, elem := range *lst {
				if reflect.TypeOf(elem).String() == "[]etree.Attr" {
					conversion := elem.([]etree.Attr)
					for _, i := range conversion {
						attr[i.Key] = i.Value
					}
				} else {
					conversion := elem.(string)
					if index == 0 && lstIndex == 0 {
						res, err := xmlProcessing(tg["XMLName"])
						if err != nil {
							return "", nil, "", errors.Annotate(err, "xmlProcessing")
						}
						elemName = res
					} else if index == 0 {
						res, err := xmlProcessing(tg[conversion])
						if err != nil {
							return "", nil, "", errors.Annotate(err, "xmlProcessing")
						}
						elemName = res
					} else {
						value = conversion
					}
				}
			}
		}
	}
	return elemName, attr, value, nil
}

func xmlProcessing(tg string) (string, error) {
	r, _ := regexp.Compile(`"(.*?)"`)
	str := r.FindStringSubmatch(tg)
	if len(str) == 0 {
		return "", errors.New("out of range")
	}
	attr := strings.Index(str[1], ",attr")
	omit := strings.Index(str[1], ",omitempty")
	attrOmit := strings.Index(str[1], ",attr,omitempty")
	omitAttr := strings.Index(str[1], ",omitempty,attr")

	if attr > -1 && attrOmit == -1 && omitAttr == -1 {
		return str[1][0:attr], nil
	} else if omit > -1 && attrOmit == -1 && omitAttr == -1 {
		return str[1][0:omit], nil
	} else if attr == -1 && omit == -1 {
		return str[1], nil
	} else if attrOmit > -1 {
		return str[1][0:attrOmit], nil
	} else {
		return str[1][0:omitAttr], nil
	}
}

func mapProcessing(mapVar []map[string]string) []map[string]string {
	for indx := 0; indx < len(mapVar); indx++ {
		element := mapVar[indx]
		for _, value := range element {
			if value == "" {
				mapVar = append(mapVar[:indx], mapVar[indx+1:]...)
				indx--
			}
			if strings.Index(value, ",attr") != -1 {
				mapVar = append(mapVar[:indx], mapVar[indx+1:]...)
				indx--
			}
		}
	}
	return mapVar
}

func soapHandling(tp interface{}, tags *[]map[string]string) {
	s := reflect.ValueOf(tp).Elem()
	typeOfT := s.Type()
	if s.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tmp, ok := typeOfT.FieldByName(typeOfT.Field(i).Name)
		if !ok {
			Logger.Debug().Str("field", typeOfT.Field(i).Name).Msg("reflection failed")
		}
		*tags = append(*tags, map[string]string{typeOfT.Field(i).Name: string(tmp.Tag)})
		subStruct := reflect.New(reflect.TypeOf(f.Interface()))
		soapHandling(subStruct.Interface(), tags)
	}
}

func xmlUnmarshal(elems []*etree.Element, data *[][]interface{}, mas *[]string) {
	for _, elem := range elems {
		*data = append(*data, []interface{}{elem.Tag, elem.Attr, elem.Text()})
		*mas = append(*mas, "Push")
		xmlUnmarshal(elem.FindElements("./*"), data, mas)
		*mas = append(*mas, "Pop")
	}
}

func ident(mas *[]string) {
	var buffer string
	for _, j := range *mas {
		buffer += j + " "
	}
	buffer = strings.Replace(buffer, "Push Pop ", "PushPop ", -1)
	buffer = strings.TrimSpace(buffer)
	*mas = strings.Split(buffer, " ")
}

func main() {
	RunApi() // Start the API
}
