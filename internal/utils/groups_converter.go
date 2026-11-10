// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/legacy_converter.go
package utils

import (
	"strings"

	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertLegacyServiceGroup(legacyGroup legacy.ServiceGroup) models.ServiceGroup {
	groupType := getLegacyGroupType(legacyGroup)
	serverType := getStringFromPayload(legacyGroup.Payload, "mode")

		if groupType == "server" && serverType == "" {
		serverType = getLegacyServerSubtype(legacyGroup)
	}

	group := models.ServiceGroup{
		ID:             legacyGroup.ID,
		Name:           legacyGroup.Name,
		Type:           legacyGroup.TypeObject.Slug,
		GroupType:      groupType,
		ServiceHandler: legacyGroup.ServiceHandler,
		Description:    legacyGroup.Description,
		Location:       getStringFromPayload(legacyGroup.Payload, "label"),
		CountryCode:    getStringFromPayload(legacyGroup.Payload, "code"),
		ServerType:     serverType,
		IsDisabled:     getLegacyBoolFromPayload(legacyGroup.Payload, "isDisabled"),
		Features:       getLegacyFeatures(legacyGroup.LocaledPayload),
	}

		if groupType == "server" {
		features := getLegacyFeatures(legacyGroup.LocaledPayload)
		group.CPUModel = extractCPUModel(features)
		group.CPUFrequency = extractCPUFrequency(features)
		group.NetworkSpeed = extractNetworkSpeed(features)
		group.IPv4Count = extractIPv4Count(features)
		group.IPv6Subnet = extractIPv6Subnet(features)
	}

	return group
}

func getLegacyGroupType(legacyGroup legacy.ServiceGroup) string {
	mode := getStringFromPayload(legacyGroup.Payload, "mode")
	code := getStringFromPayload(legacyGroup.Payload, "code")
	role := getLegacyRole(legacyGroup.Role)

		if code != "" {
		return "location"
	}

		if mode != "" {
		return "server"
	}

	// 3. Special services
	if isSpecialService(legacyGroup.TypeObject.Slug) {
		return "special"
	}

		if role == "location" {
		return "geography"
	}

	return "unknown"
}

func getLegacyRole(role interface{}) string {
	if role == nil {
		return ""
	}
	if str, ok := role.(string); ok {
		return str
	}
	return ""
}

func getLegacyFeatures(localed map[string]interface{}) string {
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

		if pretty, exists := localed["prettyLocaledPayload"]; exists {
		if prettyMap, ok := pretty.(map[string]interface{}); ok {
			if features, exists := prettyMap["features"]; exists {
				if str, ok := features.(string); ok {
					return str
				}
			}
		}
	}

	return ""
}

func getLegacyServerSubtype(legacyGroup legacy.ServiceGroup) string {
	mode := getStringFromPayload(legacyGroup.Payload, "mode")
	label := getStringFromPayload(legacyGroup.Payload, "label")

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

func getLegacyBoolFromPayload(payload map[string]interface{}, key string) bool {
	if val, exists := payload[key]; exists {
		switch v := val.(type) {
		case bool:
			return v
		case string:
						return v == "true" || v == "1"
		}
	}
	return false
}

func ConvertLegacyServiceGroups(legacyGroups []legacy.ServiceGroup) []models.ServiceGroup {
	result := make([]models.ServiceGroup, len(legacyGroups))
	for i, legacyGroup := range legacyGroups {
		result[i] = ConvertLegacyServiceGroup(legacyGroup)
	}
	return result
}
