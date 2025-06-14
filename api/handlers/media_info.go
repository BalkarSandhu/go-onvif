package handlers

import (
	"go-onvif/api/extractors"
	"go-onvif/api/service"
	onvif "go-onvif/internal"
)

func GetMediaInfo(ip string, dev *onvif.Device) interface{} {
	profilesResp, err := service.CallMediaMethod("GetProfiles", dev, nil)
	if err != nil {
		return extractors.ErrorResponse{Error: err.Error()}
	}

	profiles, err := extractors.ParseProfilesResponse(profilesResp, ip)
	if err != nil {
		return extractors.ErrorResponse{Error: err.Error()}
	}

	importantProfiles := make([]extractors.ProfileInfo, 0, len(profiles))
	for _, profile := range profiles {
		profileInfo := extractors.ExtractProfileInfo(dev, profile)
		importantProfiles = append(importantProfiles, profileInfo)
	}

	return importantProfiles
}
