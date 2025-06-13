package api

func (s *APIServer) SetupRoutes() {
	s.router.Use(s.loggerMiddleware())
	s.router.Use(s.rateLimitMiddleware())

	s.router.POST("/:service/:method", s.handleServiceMethod)
	s.router.GET("/discovery", s.handleDiscovery)
	s.router.GET("/scan", s.handleDeviceScan)
	s.router.GET("/device_cache", s.getDeviceCache)
	s.router.GET("/device_details", s.getDeviceDetails)
}
