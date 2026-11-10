// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/service_create_utils.go
package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertToLegacy_ServiceCreateRequest(req models.ServiceCreateRequest) legacy.ServiceCreateRequest {
	return legacy.ServiceCreateRequest{
		Method:      "balance",
		Count:       1,
		Term:        req.PaymentTerm,
		Name:        req.Name,
		ProductID:   req.ProductID,
		AutoProlong: req.AutoProlong,
		Backups:     false,
		Parameters: legacy.ServiceCreateParameters{
			OS:     req.OS,
			Recipe: stringToPtr(req.Recipe),
			IsoURL: req.IsoURL,
		},
	}
}

func ConvertFromLegacy_ServiceCreateResponse(resp legacy.ServiceCreateResponse) models.ServiceCreateResponse {
		if len(resp.Data.Items) == 0 {
		return models.ServiceCreateResponse{
			ID:     0,
			Status: "error",
			Date:   FormatDateFromUnix(0),
		}
	}

	item := resp.Data.Items[0]
	transaction := resp.Data.Transaction

		var serviceID int64
	if len(item.CreatedServiceIds) > 0 {
		serviceID = item.CreatedServiceIds[0]
	}

		locationName := extractLocationName(item.Product)

		price := FormatPrice(item.IndividualPrice)

		transactionAmountValue := transaction.Amount
	if transactionAmountValue < 0 {
		transactionAmountValue = -transactionAmountValue
	}
	transactionAmount := FormatPrice(transactionAmountValue)

	return models.ServiceCreateResponse{
		ID:                serviceID,
		OrderID:           item.ID,
		Status:            item.Status,
		Date:              FormatDateFromUnix(item.CreatedAt),
		ProductId:         item.ProductID,
		ProductType:       item.Product.Type,
		GroupId:           item.Product.GroupID,
		ProductName:       item.Product.Name,
		LocationName:      locationName,
		Term:              item.Term,
		Price:             price,
		TransactionAmount: transactionAmount,
	}
}

func extractLocationName(product legacy.Product) string {
	if product.Group.Payload != nil {
		if label, exists := product.Group.Payload["label"]; exists {
			if labelStr, ok := label.(string); ok {
				return labelStr
			}
		}
	}
	return ""
}
