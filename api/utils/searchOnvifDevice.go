package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func IsONVIFDevice(ip string) bool {
	soapPayload := `
<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope">
    <Header/>
    <Body>
        <GetServices xmlns="http://www.onvif.org/ver10/device/wsdl">
            <IncludeCapability>false</IncludeCapability>
        </GetServices>
    </Body>
</Envelope>`

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post("http://"+ip+"/onvif/device_service", "application/soap+xml", bytes.NewBuffer([]byte(soapPayload)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	match, _ := regexp.Match(`(?i)<.*Envelope.*>`, body)
	return match
}
