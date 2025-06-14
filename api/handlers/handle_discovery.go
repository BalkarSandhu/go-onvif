package handlers

import (
	"net/http"
	"regexp"
	"strings"

	wsdiscovery "go-onvif/internal/ws-discovery"

	"github.com/beevik/etree"
	"github.com/gin-gonic/gin"
)

func HandleDiscovery(c *gin.Context) {
	iface := c.Query("interface")

	devices, err := wsdiscovery.SendProbe(iface, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Discovery failed: " + err.Error()})
		return
	}

	var discovered []map[string]string

	for _, deviceXML := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(deviceXML); err != nil {
			continue
		}

		xaddrEl := doc.Root().FindElement("./Body/ProbeMatches/ProbeMatch/XAddrs")
		scopeEl := doc.Root().FindElement("./Body/ProbeMatches/ProbeMatch/Scopes")

		if xaddrEl == nil {
			continue
		}

		xaddr := strings.Split(xaddrEl.Text(), " ")[0]
		url := strings.Split(xaddr, "/")[2]

		found := false
		for _, d := range discovered {
			if d["url"] == url {
				found = true
				break
			}
		}
		if found {
			continue
		}

		name := "Unknown"
		if scopeEl != nil {
			if m := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/([A-Za-z0-9-]+)`).FindStringSubmatch(scopeEl.Text()); len(m) > 1 {
				name = m[1]
			}
		}

		discovered = append(discovered, map[string]string{"url": url, "name": name})
	}

	c.JSON(http.StatusOK, discovered)
}
