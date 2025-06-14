package extractors

import (
	"encoding/json"
	"fmt"
)

func ParseProfilesResponse(profilesResp interface{}, ip string) ([]map[string]interface{}, error) {
	jsonData, err := json.Marshal(profilesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profiles response: %w", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profiles JSON: %w", err)
	}

	profilesRaw, exists := parsed["Profiles"]
	if !exists {
		return nil, fmt.Errorf("profiles not found in response")
	}

	profilesSlice, ok := profilesRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("profiles format invalid - expected array")
	}

	profiles := make([]map[string]interface{}, 0, len(profilesSlice))
	for _, profile := range profilesSlice {
		if profileMap, ok := profile.(map[string]interface{}); ok {
			profiles = append(profiles, profileMap)
		}
	}

	return profiles, nil
}
