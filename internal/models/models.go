// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/models/models.go
package models

import (
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/next"
)

type TerraformService struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	IP            string   `json:"ip"`
	IPv4Addresses []string `json:"ipv4_addresses,omitempty"`
	IPv6Addresses []string `json:"ipv6_addresses,omitempty"`
	Price         float64  `json:"price"`
	PaymentTerm   string   `json:"paymentTerm"`
	AutoProlong   bool     `json:"autoProlong"`
	CreatedAt     string   `json:"createdAt"`
	ExpiresAt     string   `json:"expiresAt"`
	Status        string   `json:"status"`
	TypeSlug      string   `json:"typeSlug"`
	ProductName   string   `json:"productName"`
	LocationCode  string   `json:"locationCode"`
}

type TerraformServiceDetailed struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	TypeSlug     string `json:"type_slug"`
	LocationCode string `json:"location_code"`

	IP   string   `json:"ip"`
	IPs  []string `json:"ips,omitempty"`
	IPv6 string   `json:"ipv6,omitempty"`

	Price        float64 `json:"price"`
	PriceDisplay string  `json:"price_display"`
	PaymentTerm  string  `json:"payment_term"`
	AutoProlong  bool    `json:"auto_prolong"`

	CPU     int `json:"cpu,omitempty"`
	RAM     int `json:"ram,omitempty"`
	Storage int `json:"storage,omitempty"`
	IPCount int `json:"ip_count,omitempty"`

	ProductName string `json:"product_name"`
	OS          string `json:"os,omitempty"`
	Username    string `json:"username,omitempty"`

	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
	CreatedDate string `json:"created_date"`
	ExpiresDate string `json:"expires_date"`

	Capabilities []string               `json:"capabilities,omitempty"`
	Backups      bool                   `json:"backups"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
}

type Product struct {
	ID                   int64                           `json:"id"`
	Type                 string                          `json:"type"`
	GroupID              int64                           `json:"groupId"`
	Name                 string                          `json:"name"`
	Configuration        []ConfigurationItem             `json:"configuration"`
	Prices               map[string]float64              `json:"prices"`
	SummaryConfiguration map[string]ConfigurationSummary `json:"summaryConfiguration"`
	ServiceHandler       string                          `json:"serviceHandler"`
}

type ConfigurationItem struct {
	Max    float64                `json:"max"`
	Base   float64                `json:"base"`
	Slug   string                 `json:"slug"`
	Type   string                 `json:"type"`
	Prices map[string]interface{} `json:"prices"`
}

type ConfigurationSummary struct {
	Max  float64 `json:"max"`
	Base float64 `json:"base"`
	Slug string  `json:"slug"`
	Type string  `json:"type"`
}

type ServiceType struct {
	Slug    string                 `json:"slug"`
	Order   int                    `json:"order"`
	Names   map[string]string      `json:"names"`
	Payload map[string]interface{} `json:"payload"`
	Name    string                 `json:"name"`
}

type OperatingSystem struct {
	next.OperatingSystem
}

type Recipe struct {
	ID              int            `json:"id"`
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Description     string         `json:"description"`
	Tags            []string       `json:"tags"`
	Enabled         bool           `json:"enabled"`
	OS              int            `json:"os"`
	OperatingSystem string         `json:"operating_system"`
	Targets         map[string]int `json:"targets"`
}
