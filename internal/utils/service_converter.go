// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/service_converter.go
package utils

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/legacy"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/next"
)

func ConvertNextService(nextService next.Service) models.TerraformService {
	var ipv4Addresses, ipv6Addresses []string
	if isIPv6Address(nextService.IP) {
		ipv6Addresses = compactUniqueStrings([]string{nextService.IP})
	} else {
		ipv4Addresses = compactUniqueStrings([]string{nextService.IP})
	}

	return models.TerraformService{
		ID:            nextService.ID,
		Name:          nextService.Name,
		IP:            nextService.IP,
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		Price:         nextService.Price,
		PaymentTerm:   nextService.PaymentTerm,
		AutoProlong:   nextService.AutoProlong,
		CreatedAt:     nextService.CreatedAt,
		ExpiresAt:     nextService.ExpiresAt,
		Status:        nextService.Status,
		TypeSlug:      nextService.TypeSlug,
		ProductName:   nextService.ProductName,
		LocationCode:  nextService.LocationCode,
	}
}

func ConvertLegacyService(legacyService legacy.ServiceVPS) models.TerraformService {
	ipv4Addresses := extractIPv4Addresses(legacyService.IP, legacyService.IPs)
	ipv6Addresses := extractIPv6Addresses(legacyService.IP, legacyService.IPv6)

	tfService := models.TerraformService{
		ID:            legacyService.ID,
		Name:          legacyService.Name,
		IP:            legacyService.IP,
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		PaymentTerm:   legacyService.PaymentTerm,
		AutoProlong:   legacyService.AutoProlong,
		Status:        legacyService.Status,
		LocationCode:  legacyService.LocationCode,
	}

	if legacyService.Product.Name != "" {
		tfService.ProductName = legacyService.Product.Name
	} else {
		tfService.ProductName = "Unknown"
	}

	if legacyService.Product.Type != "" {
		tfService.TypeSlug = legacyService.Product.Type
	} else {
		tfService.TypeSlug = "unknown"
	}

	tfService.CreatedAt = convertUnixToISO(int64(legacyService.Timestamps.CreatedAt))
	tfService.ExpiresAt = convertUnixToISO(int64(legacyService.Timestamps.ExpiresAt))

	tfService.Price = calculatePriceFromLegacy(legacyService)

	return tfService
}

func calculatePriceFromLegacy(legacyService legacy.ServiceVPS) float64 {
	if legacyService.IndividualPrices != nil {
		if price, exists := legacyService.IndividualPrices[legacyService.PaymentTerm]; exists {
			return float64(price)
		}
	}

	if legacyService.RawPrices != nil {
		if price, exists := legacyService.RawPrices[legacyService.PaymentTerm]; exists {
			return float64(price)
		}
	}

	if legacyService.Product.Prices.Hour > 0 || legacyService.Product.Prices.Month > 0 {
		switch legacyService.PaymentTerm {
		case "hour":
			return legacyService.Product.Prices.Hour
		case "month":
			return legacyService.Product.Prices.Month
		case "year":
			return legacyService.Product.Prices.Year
		case "half_year":
			return legacyService.Product.Prices.HalfYear
		case "quarter_year":
			return legacyService.Product.Prices.QuarterYear
		}
	}

	if legacyService.Product.RawPrices != (legacy.ProductPrices{}) {
		switch legacyService.PaymentTerm {
		case "hour":
			if legacyService.Product.RawPrices.Hour > 0 {
				return legacyService.Product.RawPrices.Hour
			}
		case "month":
			if legacyService.Product.RawPrices.Month > 0 {
				return legacyService.Product.RawPrices.Month
			}
		case "year":
			if legacyService.Product.RawPrices.Year > 0 {
				return legacyService.Product.RawPrices.Year
			}
		case "half_year":
			if legacyService.Product.RawPrices.HalfYear > 0 {
				return legacyService.Product.RawPrices.HalfYear
			}
		case "quarter_year":
			if legacyService.Product.RawPrices.QuarterYear > 0 {
				return legacyService.Product.RawPrices.QuarterYear
			}
		}
	}

	if legacyService.Product.IndividualPrices != (legacy.ProductPrices{}) {
		switch legacyService.PaymentTerm {
		case "hour":
			if legacyService.Product.IndividualPrices.Hour > 0 {
				return legacyService.Product.IndividualPrices.Hour
			}
		case "month":
			if legacyService.Product.IndividualPrices.Month > 0 {
				return legacyService.Product.IndividualPrices.Month
			}
		case "year":
			if legacyService.Product.IndividualPrices.Year > 0 {
				return legacyService.Product.IndividualPrices.Year
			}
		case "half_year":
			if legacyService.Product.IndividualPrices.HalfYear > 0 {
				return legacyService.Product.IndividualPrices.HalfYear
			}
		case "quarter_year":
			if legacyService.Product.IndividualPrices.QuarterYear > 0 {
				return legacyService.Product.IndividualPrices.QuarterYear
			}
		}
	}

	return 0
}

func ConvertNextServices(nextServices []next.Service) []models.TerraformService {
	var result []models.TerraformService
	for _, service := range nextServices {
		result = append(result, ConvertNextService(service))
	}
	return result
}

func ConvertLegacyServices(legacyServices []legacy.ServiceVPS) []models.TerraformService {
	var result []models.TerraformService
	for _, service := range legacyServices {
		result = append(result, ConvertLegacyService(service))
	}
	return result
}

func ConvertLegacyVPSToTerraform(vpsService legacy.ServiceVPS) models.TerraformServiceDetailed {
	tfService := models.TerraformServiceDetailed{
		ID:           vpsService.ID,
		Name:         vpsService.Name,
		Status:       vpsService.Status,
		LocationCode: vpsService.LocationCode,
		AutoProlong:  vpsService.AutoProlong,
		Backups:      vpsService.Backups,
		PaymentTerm:  vpsService.PaymentTerm,
		Capabilities: vpsService.Capabilities,
		Payload:      vpsService.Payload,

		IP:   vpsService.IP,
		IPs:  extractIPs(vpsService.IPs),
		IPv6: extractIPv6(vpsService.IPv6),

		ProductName: vpsService.Product.Name,
		TypeSlug:    getTypeSlug(vpsService.Product),

		OS:       getOSDisplayName(vpsService.Parameters.OS),
		Username: vpsService.Parameters.Username,
	}

	tfService.Price = getPriceInEuros(vpsService.RawPrices, vpsService.PaymentTerm)
	tfService.PriceDisplay = formatPriceDisplay(tfService.Price, vpsService.PaymentTerm)

	if vpsService.SummaryConfiguration != nil {
		tfService.CPU = getCPUCount(vpsService.SummaryConfiguration["cpu"])
		tfService.RAM = getRAMCount(vpsService.SummaryConfiguration["ram"])
		tfService.Storage = getStorageCount(vpsService.SummaryConfiguration["rom"])
		tfService.IPCount = getIPCount(vpsService.SummaryConfiguration["ip"])
	}

	tfService.CreatedAt = convertUnixToISO(int64(vpsService.Timestamps.CreatedAt))
	tfService.ExpiresAt = convertUnixToISO(int64(vpsService.Timestamps.ExpiresAt))
	tfService.CreatedDate = formatDisplayDate(int64(vpsService.Timestamps.CreatedAt))
	tfService.ExpiresDate = formatDisplayDate(int64(vpsService.Timestamps.ExpiresAt))

	return tfService
}

func ConvertNextToTerraform(nextService next.Service) models.TerraformServiceDetailed {
	return models.TerraformServiceDetailed{
		ID:           nextService.ID,
		Name:         nextService.Name,
		Status:       nextService.Status,
		LocationCode: nextService.LocationCode,
		AutoProlong:  nextService.AutoProlong,
		PaymentTerm:  nextService.PaymentTerm,
		Capabilities: nextService.Capabilities,
		Payload:      nextService.Payload,

		IP: nextService.IP,

		ProductName: nextService.ProductName,
		TypeSlug:    nextService.TypeSlug,

		Price: float64(nextService.Price) / 100, PriceDisplay: formatPriceDisplay(float64(nextService.Price)/100, nextService.PaymentTerm),

		CreatedAt:   nextService.CreatedAt,
		ExpiresAt:   nextService.ExpiresAt,
		CreatedDate: formatISOToDisplayDate(nextService.CreatedAt),
		ExpiresDate: formatISOToDisplayDate(nextService.ExpiresAt),
	}
}

func getPriceInEuros(rawPrices map[string]int, paymentTerm string) float64 {
	if rawPrices == nil {
		return 0
	}

	priceCents := rawPrices[paymentTerm]
	if priceCents == 0 && paymentTerm != "month" {
		priceCents = rawPrices["month"]
	}

	return float64(priceCents) / 100
}

func formatPriceDisplay(price float64, paymentTerm string) string {
	periods := map[string]string{
		"hour":         "Час",
		"month":        "Месяц",
		"year":         "Год",
		"half_year":    "Полгода",
		"quarter_year": "Квартал",
	}

	period := periods[paymentTerm]
	if period == "" {
		period = paymentTerm
	}

	return fmt.Sprintf("%.2f € / %s", price, period)
}

func extractIPs(ips []legacy.IPAddress) []string {
	var result []string
	for _, ip := range ips {
		result = append(result, ip.Value)
	}
	return result
}

func extractIPv4Addresses(primary string, ips []legacy.IPAddress) []string {
	var values []string
	if !isIPv6Address(primary) {
		values = append(values, primary)
	}
	for _, ip := range ips {
		values = append(values, ip.Value)
	}
	return compactUniqueStrings(values)
}

func isIPv6Address(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if parsed := net.ParseIP(value); parsed != nil {
		return parsed.To4() == nil
	}
	if parsed, _, err := net.ParseCIDR(value); err == nil {
		return parsed.To4() == nil
	}
	return strings.Contains(value, ":")
}

func extractIPv6Addresses(primary string, ipv6 []legacy.IPv6Address) []string {
	var values []string
	if isIPv6Address(primary) {
		values = append(values, primary)
	}
	for _, subnet := range ipv6 {
		values = append(values, subnet.Value)
		for _, raw := range subnet.IPs {
			switch value := raw.(type) {
			case string:
				values = append(values, value)
			case map[string]interface{}:
				for _, key := range []string{"value", "address", "ip"} {
					if address, ok := value[key].(string); ok {
						values = append(values, address)
						break
					}
				}
			}
		}
	}
	return compactUniqueStrings(values)
}

func compactUniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func extractIPv6(ipv6 []legacy.IPv6Address) string {
	if len(ipv6) > 0 {
		return ipv6[0].Value
	}
	return ""
}

func getTypeSlug(product legacy.Product) string {
	if product.Type != "" {
		return product.Type
	}
	if product.TypeObject.Slug != "" {
		return product.TypeObject.Slug
	}
	return "unknown"
}

func getCPUCount(config legacy.ConfigurationItem) int {
	return config.Count
}

func getRAMCount(config legacy.ConfigurationItem) int {
	return config.Count
}

func getStorageCount(config legacy.ConfigurationItem) int {
	return config.Count
}

func getIPCount(config legacy.ConfigurationItem) int {
	return config.Count
}

// OS mapping
func getOSDisplayName(osSlug string) string {
	mapping := map[string]string{
		"ubuntu_2404":  "Ubuntu 24.04",
		"ubuntu_2204":  "Ubuntu 22.04",
		"ubuntu_2004":  "Ubuntu 20.04",
		"centos_7":     "CentOS 7",
		"debian_11":    "Debian 11",
		"windows_2019": "Windows Server 2019",
	}

	if name, exists := mapping[osSlug]; exists {
		return name
	}
	return osSlug
}

func convertUnixToISO(unixTimestamp int64) string {
	if unixTimestamp == 0 {
		return ""
	}
	return time.Unix(unixTimestamp, 0).UTC().Format(time.RFC3339)
}

func formatDisplayDate(unixTimestamp int64) string {
	if unixTimestamp == 0 {
		return ""
	}

	t := time.Unix(unixTimestamp, 0)
	monthNames := map[time.Month]string{
		1: "января", 2: "февраля", 3: "марта", 4: "апреля",
		5: "мая", 6: "июня", 7: "июля", 8: "августа",
		9: "сентября", 10: "октября", 11: "ноября", 12: "декабря",
	}

	return fmt.Sprintf("%d %s %d г.", t.Day(), monthNames[t.Month()], t.Year())
}

func formatISOToDisplayDate(isoDate string) string {
	if isoDate == "" {
		return ""
	}

	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}

	monthNames := map[time.Month]string{
		1: "января", 2: "февраля", 3: "марта", 4: "апреля",
		5: "мая", 6: "июня", 7: "июля", 8: "августа",
		9: "сентября", 10: "октября", 11: "ноября", 12: "декабря",
	}

	return fmt.Sprintf("%d %s %d г.", t.Day(), monthNames[t.Month()], t.Year())
}
