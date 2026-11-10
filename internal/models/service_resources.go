// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/models/service_resources.go
package models

type Service struct {
	ID               int64                  `json:"id"`
	OwnerID          int64                  `json:"ownerId"`
	ProductID        int64                  `json:"productId"`
	Name             string                 `json:"name"`
	IP               string                 `json:"ip"`
	IPv4Addresses    []string               `json:"ipv4_addresses,omitempty"`
	IPv6Addresses    []string               `json:"ipv6_addresses,omitempty"`
	PaymentTerm      string                 `json:"paymentTerm"`
	Parameters       map[string]interface{} `json:"parameters"`
	SecureParameters map[string]interface{} `json:"secureParameters"`
	AutoProlong      bool                   `json:"autoProlong"`
	Backups          bool                   `json:"backups"`
	Status           string                 `json:"status"`
	LastStatus       string                 `json:"lastStatus"`
	Product          Product                `json:"product"`
	LocationCode     string                 `json:"locationCode"`
	Prices           map[string]Price       `json:"prices"`
	CurrentStatus    string                 `json:"currentStatus"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
	ExpiresAt        string                 `json:"expiresAt"`
	PurchasedAt      string                 `json:"purchasedAt"`
}

type Price struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type ServiceCreateRequest struct {
	Name        string `json:"name"`
	ProductID   int64  `json:"product_id"`
	PaymentTerm string `json:"payment_term"`
	AutoProlong bool   `json:"auto_prolong"`
	OS          string `json:"os,omitempty"`
	Recipe      string `json:"recipe,omitempty"`
	IsoURL      string `json:"iso_url,omitempty"`
}

type ServiceCreateResponse struct {
	ID      int64  `json:"id"`
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
	Date    string `json:"date"`

	ProductId   int64  `json:"product_id"`
	ProductType string `json:"product_type"`
	GroupId     *int64 `json:"group_id"`
	ProductName string `json:"product_name"`

	LocationName string `json:"location_name"`

	Term              string `json:"term"`
	Price             string `json:"price"`
	TransactionAmount string `json:"transaction_amount"`
}

type ServiceUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	AutoProlong *bool   `json:"auto_prolong,omitempty"`
}
