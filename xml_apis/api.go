package api

import (
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"

	"github.com/BalkarSandhu/go-onvif/gosoap"
	"github.com/BalkarSandhu/go-onvif/networking"
	"github.com/BalkarSandhu/go-onvif/onvif"
	wsdiscovery "github.com/BalkarSandhu/go-onvif/ws-discovery"
)

// Config holds configuration for the API server
type Config struct {
	Port           string
	LogLevel       string
	RateLimitReqs  int
	RateLimitBurst int
}

// APIServer encapsulates the HTTP server and its dependencies
type APIServer struct {
	config  Config
	router  *gin.Engine
	logger  zerolog.Logger
	limiter *rate.Limiter
}

// NewAPIServer creates a new instance of APIServer
func NewAPIServer(config Config) *APIServer {
	// Configure logging
	logLevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger().Level(logLevel)

	// Configure rate limiter
	limiter := rate.NewLimiter(rate.Limit(config.RateLimitReqs), config.RateLimitBurst)

	// Create gin router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add custom logger middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Msg("Request")
	})

	return &APIServer{
		config:  config,
		router:  router,
		logger:  logger,
		limiter: limiter,
	}
}

// rateLimitMiddleware applies rate limiting to requests
func (s *APIServer) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetupRoutes configures all the routes for the API
func (s *APIServer) SetupRoutes() {
	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "username", "password", "xaddr", "interface")
	s.router.Use(cors.New(corsConfig), s.rateLimitMiddleware())

	// Service endpoint
	s.router.POST("/:service/:method", s.handleServiceMethod)

	// Discovery endpoint
	s.router.GET("/discovery", s.handleDiscovery)
}

// Run starts the HTTP server
func (s *APIServer) Run() error {
	return s.router.Run(":" + s.config.Port)
}

// handleServiceMethod handles POST requests to service methods
func (s *APIServer) handleServiceMethod(c *gin.Context) {
	serviceName := c.Param("service")
	methodName := c.Param("method")
	username := c.GetHeader("username")
	pass := c.GetHeader("password")
	xaddr := c.GetHeader("xaddr")

	acceptedData, err := c.GetRawData()
	if err != nil {
		s.logger.Debug().Err(err).Msg("Failed to get raw data")
		c.XML(http.StatusBadRequest, err.Error())
		return
	}

	message, err := s.callNecessaryMethod(serviceName, methodName, string(acceptedData), username, pass, xaddr)
	if err != nil {
		c.XML(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, message)
}

// handleDiscovery handles GET requests for device discovery
func (s *APIServer) handleDiscovery(c *gin.Context) {
	interfaceName := c.GetHeader("interface")

	devices, err := wsdiscovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		s.logger.Error().Err(err).Msg("Error during device discovery")
		c.String(http.StatusInternalServerError, "error")
		return
	}

	response := "["

	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			s.logger.Error().Err(err).Msg("Error parsing discovery response")
			continue
		}

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
			continue
		}

		for _, scope := range scopes {
			re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
			match := re.FindStringSubmatch(scope.Text())
			if len(match) > 0 {
				response += `"name":"` + path.Base(match[0]) + `"`
			}
		}

		response += "},"
	}

	response = strings.TrimRight(response, ",")
	response += "]"
	c.String(http.StatusOK, response)
}

// callNecessaryMethod calls the appropriate ONVIF service method
func (s *APIServer) callNecessaryMethod(serviceName, methodName, acceptedData, username, password, xaddr string) (string, error) {
	var methodStruct interface{}
	var err error

	switch strings.ToLower(serviceName) {
	case "device":
		methodStruct, err = getDeviceStructByName(methodName)
	case "ptz":
		methodStruct, err = getPTZStructByName(methodName)
	case "media":
		methodStruct, err = getMediaStructByName(methodName)
	default:
		return "", errors.New("there is no such service")
	}

	if err != nil {
		return "", errors.Annotate(err, "getStructByName")
	}

	resp, err := xmlAnalize(methodStruct, &acceptedData)
	if err != nil {
		return "", errors.Annotate(err, "xmlAnalize")
	}

	endpoint, err := getEndpoint(serviceName, xaddr)
	if err != nil {
		return "", errors.Annotate(err, "getEndpoint")
	}

	soap := gosoap.NewEmptySOAP()
	soap.AddStringBodyContent(*resp)
	soap.AddRootNamespaces(onvif.Xlmns)
	soap.AddWSSecurity(username, password)

	servResp, err := networking.SendSoap(new(http.Client), endpoint, soap.String())
	if err != nil {
		return "", errors.Annotate(err, "SendSoap")
	}

	rsp, err := io.ReadAll(servResp.Body)
	if err != nil {
		return "", errors.Annotate(err, "ReadAll")
	}

	defer servResp.Body.Close()

	return string(rsp), nil
}

// getEndpoint returns the endpoint URL for a service
func getEndpoint(service, xaddr string) (string, error) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: xaddr})
	if err != nil {
		return "", errors.Annotate(err, "NewDevice")
	}

	pkg := strings.ToLower(service)

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

// xmlAnalize processes XML data for SOAP requests
func xmlAnalize(methodStruct interface{}, acceptedData *string) (*string, error) {
	test := make([]map[string]string, 0)      // tags
	testunMarshal := make([][]interface{}, 0) // data
	var mas []string                          // idnt

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

		if mas[lstIndex] == "Push" && lstIndex == 0 {
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

// xmlMaker creates XML elements from struct data
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

// xmlProcessing extracts XML tag names from struct tags
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

// mapProcessing filters map data
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

// soapHandling creates tags from struct reflection
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
			// Use server logger
			// Logger.Debug().Str("field", typeOfT.Field(i).Name).Msg("reflection failed")
			continue
		}

		*tags = append(*tags, map[string]string{typeOfT.Field(i).Name: string(tmp.Tag)})
		subStruct := reflect.New(reflect.TypeOf(f.Interface()))
		soapHandling(subStruct.Interface(), tags)
	}
}

// xmlUnmarshal converts etree elements to data structures
func xmlUnmarshal(elems []*etree.Element, data *[][]interface{}, mas *[]string) {
	for _, elem := range elems {
		*data = append(*data, []interface{}{elem.Tag, elem.Attr, elem.Text()})
		*mas = append(*mas, "Push")
		xmlUnmarshal(elem.FindElements("./*"), data, mas)
		*mas = append(*mas, "Pop")
	}
}

// ident processes identification markers
func ident(mas *[]string) {
	var buffer string
	for _, j := range *mas {
		buffer += j + " "
	}

	buffer = strings.Replace(buffer, "Push Pop ", "PushPop ", -1)
	buffer = strings.TrimSpace(buffer)
	*mas = strings.Split(buffer, " ")
}
