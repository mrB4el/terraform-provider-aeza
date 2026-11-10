// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/service_get_utils.go
package utils

import (
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/legacy"
)

func ConvertLegacyServiceGetToTerraform(legacyService legacy.ServiceGet) models.Service {
	ipv4Addresses := extractIPv4Addresses(legacyService.IP, legacyService.IPs)
	ipv6Addresses := extractIPv6Addresses(legacyService.IP, legacyService.IPv6)

	ip := legacyService.IP
	if ip == "" && len(ipv4Addresses) > 0 {
		ip = ipv4Addresses[0]
	}

	parameters := make(map[string]interface{})

	parameters["os"] = legacyService.Parameters.OS
	parameters["isoUrl"] = legacyService.Parameters.ISOURL
	if legacyService.Parameters.Recipe != nil {
		parameters["recipe"] = *legacyService.Parameters.Recipe
	}

	// if legacyService.Product.Type == "kubernetes" {
	//     parameters["kubernetes_version"] = legacyService.KubernetesVersion
	// }

	createdAt := FormatDateFromUnix(int64(legacyService.Timestamps.CreatedAt))
	expiresAt := FormatDateFromUnix(int64(legacyService.Timestamps.ExpiresAt))
	purchasedAt := FormatDateFromUnix(int64(legacyService.Timestamps.PurchasedAt))

	return models.Service{
		ID:            int64(legacyService.ID),
		OwnerID:       int64(legacyService.OwnerID),
		ProductID:     int64(legacyService.ProductID),
		Name:          legacyService.Name,
		IP:            ip,
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		PaymentTerm:   legacyService.PaymentTerm,
		Parameters:    parameters,
		AutoProlong:   legacyService.AutoProlong,
		Backups:       legacyService.Backups,
		Status:        legacyService.Status,
		LastStatus:    stringToEmpty(legacyService.LastStatus),
		Product:       ConvertLegacyProduct(legacyService.Product),
		LocationCode:  legacyService.LocationCode,
		CurrentStatus: legacyService.CurrentStatus,
		CreatedAt:     createdAt,
		UpdatedAt:     "",
		ExpiresAt:     expiresAt,
		PurchasedAt:   purchasedAt,
	}
}

func stringToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
