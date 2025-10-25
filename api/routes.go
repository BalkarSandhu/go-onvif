package api

import "go-onvif/api/handlers"

func (s *APIServer) SetupRoutes() {
	s.router.Use(s.loggerMiddleware())
	s.router.Use(s.rateLimitMiddleware())

	s.router.POST("/:service/:method", s.service.HandleServiceMethod)
	s.router.GET("/discovery", handlers.HandleDiscovery)
	s.router.GET("/scan", s.handleDeviceScan)
	s.router.GET("/device_cache", s.getDeviceCache)
	s.router.GET("/device_details", s.getDeviceDetails)

}
