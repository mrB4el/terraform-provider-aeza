// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/product_converter.go
package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertLegacyProduct(legacyProduct legacy.Product) models.Product {
	product := models.Product{
		ID:             int64(legacyProduct.ID),
		Type:           legacyProduct.Type,
		Name:           legacyProduct.Name,
		ServiceHandler: legacyProduct.ServiceHandler,
	}

		if legacyProduct.GroupID != nil {
		product.GroupID = int64(*legacyProduct.GroupID)
	}

		product.Prices = make(map[string]float64)
	if legacyProduct.Prices.Hour > 0 {
		product.Prices["hour"] = legacyProduct.Prices.Hour
	}
	if legacyProduct.Prices.Month > 0 {
		product.Prices["month"] = legacyProduct.Prices.Month
	}
	if legacyProduct.Prices.Year > 0 {
		product.Prices["year"] = legacyProduct.Prices.Year
	}
	if legacyProduct.Prices.HalfYear > 0 {
		product.Prices["half_year"] = legacyProduct.Prices.HalfYear
	}
	if legacyProduct.Prices.QuarterYear > 0 {
		product.Prices["quarter_year"] = legacyProduct.Prices.QuarterYear
	}

	// Configuration
	product.Configuration = make([]models.ConfigurationItem, len(legacyProduct.Configuration))
	for i, config := range legacyProduct.Configuration {
		product.Configuration[i] = models.ConfigurationItem{
			Max:    config.Max,  			Base:   config.Base, 			Slug:   config.Slug,
			Type:   config.Type,
			Prices: config.Prices,
		}
	}
	// SummaryConfiguration
	product.SummaryConfiguration = make(map[string]models.ConfigurationSummary)
	for key, summary := range legacyProduct.SummaryConfiguration {
		if summaryMap, ok := summary.(map[string]interface{}); ok {
			configSummary := models.ConfigurationSummary{
				Slug: getStringValue(summaryMap, "slug"),
				Type: getStringValue(summaryMap, "type"),
				Max:  getFloatValue(summaryMap, "max"),
				Base: getFloatValue(summaryMap, "base"),
			}
			product.SummaryConfiguration[key] = configSummary
		}
	}

	return product
}

func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case float32:
			return float64(v)
		}
	}
	return 0
}
