// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/groups_v2_converter.go
package utils

import (
	"strings"

	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/next"
)

func ConvertNextServiceGroup(apiGroup next.ServiceGroup) models.ServiceGroup {
	groupType := getDetailedGroupType(apiGroup)
	serverType := getStringFromPayload(apiGroup.Payload, "mode")

		if groupType == "server" && serverType == "" {
		serverType = getServerSubtype(apiGroup)
	}

	group := models.ServiceGroup{
		ID:             apiGroup.ID,
		Name:           apiGroup.Name,
		Type:           apiGroup.Type.Slug,
		GroupType:      groupType,
		ServiceHandler: apiGroup.Type.ServiceHandler,
		Description:    apiGroup.Description,
		Location:       getStringFromPayload(apiGroup.Payload, "label"),
		CountryCode:    getStringFromPayload(apiGroup.Payload, "code"),
		ServerType:     serverType,
		IsDisabled:     getBoolFromPayload(apiGroup.Payload, "isDisabled"),
		Features:       getFeatures(apiGroup.LocaledPayload),
	}

		if groupType == "server" {
		features := getFeatures(apiGroup.LocaledPayload)
		group.CPUModel = extractCPUModel(features)
		group.CPUFrequency = extractCPUFrequency(features)
		group.NetworkSpeed = extractNetworkSpeed(features)
		group.IPv4Count = extractIPv4Count(features)
		group.IPv6Subnet = extractIPv6Subnet(features)
	}

	return group
}

func getFeatures(localed map[string]interface{}) string {
	if localed == nil {
		return ""
	}

	
		if features, exists := localed["features"]; exists {
		switch v := features.(type) {
		case map[string]interface{}:
						if ru, exists := v["ru"]; exists {
				if str, ok := ru.(string); ok {
					return str
				}
			}
						if en, exists := v["en"]; exists {
				if str, ok := en.(string); ok {
					return str
				}
			}
		case string:
						return v
		}
	}

	return ""
}

func getStringFromPayload(payload map[string]interface{}, key string) string {
	if val, exists := payload[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBoolFromPayload(payload map[string]interface{}, key string) bool {
	if val, exists := payload[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func extractCPUModel(features string) string {
	lines := strings.Split(features, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Процессор") {
			return strings.TrimPrefix(line, "Процессор ")
		}
		if strings.Contains(line, "Processor") {
			return strings.TrimPrefix(line, "Processor ")
		}
	}
	return ""
}

func extractCPUFrequency(features string) string {
	lines := strings.Split(features, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Частота") {
			return strings.TrimPrefix(line, "Частота ")
		}
		if strings.Contains(line, "Frequency") {
			return strings.TrimPrefix(line, "Frequency ")
		}
	}
	return ""
}

func extractNetworkSpeed(features string) string {
	lines := strings.Split(features, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Интернет до") {
			return strings.TrimPrefix(line, "Интернет до ")
		}
		if strings.Contains(line, "Internet up to") {
			return strings.TrimPrefix(line, "Internet up to ")
		}
	}
	return ""
}

func extractIPv4Count(features string) int {
	lines := strings.Split(features, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "1 адрес IPv4") || strings.Contains(line, "1 IPv4 address") {
			return 1
		}
	}
	return 0
}

func extractIPv6Subnet(features string) string {
	lines := strings.Split(features, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "/48") {
			return "/48"
		}
		if strings.Contains(line, "/64") {
			return "/64"
		}
	}
	return ""
}

func ConvertNextServiceGroups(nextGroups []next.ServiceGroup) []models.ServiceGroup {
	result := make([]models.ServiceGroup, len(nextGroups))
	for i, apiGroup := range nextGroups {
		result[i] = ConvertNextServiceGroup(apiGroup)
	}
	return result
}

func getDetailedGroupType(apiGroup next.ServiceGroup) string {
	mode := getStringFromPayload(apiGroup.Payload, "mode")
	code := getStringFromPayload(apiGroup.Payload, "code")

		if code != "" {
		return "location"
	}

		if mode != "" {
		return "server"
	}

	// 4. Special services
	if isSpecialService(apiGroup.Type.Slug) {
		return "special"
	}

	return "unknown"
}

func isSpecialService(slug string) bool {
	specialServices := []string{"waf", "vpn", "s3", "soft"}
	for _, service := range specialServices {
		if slug == service {
			return true
		}
	}
	return false
}

func getServerSubtype(apiGroup next.ServiceGroup) string {
	mode := getStringFromPayload(apiGroup.Payload, "mode")
	label := getStringFromPayload(apiGroup.Payload, "label")

	if mode == "shared" {
		return "shared"
	}
	if mode == "dedicated" {
		return "dedicated"
	}

		if strings.Contains(label, "SHARED") {
		return "shared"
	}
	if strings.Contains(label, "DEDICATED") {
		return "dedicated"
	}

	return "server"
}
