// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/service_prolong_utils.go

package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertToLegacy_ServiceProlongRequest(req models.ServiceProlongRequest) legacy.ServiceProlongRequest {
	return legacy.ServiceProlongRequest{
		Method: req.Method,
		Term:   req.Term,
		Count:  req.Count,
	}
}

func ConvertFromLegacy_ServiceProlongResponse(legacyResp *legacy.ServiceProlongResponse) *models.ServiceProlongResponse {
	if legacyResp == nil {
		return nil
	}

	legacyTx := legacyResp.Data.Transaction

		transactionAmountValue := legacyTx.Amount
	if transactionAmountValue < 0 {
		transactionAmountValue = -transactionAmountValue
	}
	transactionAmount := FormatPrice(transactionAmountValue)

	tx := &models.ProlongedTransaction{
		ID:        legacyTx.ID,
		Amount:    transactionAmount,
		Status:    legacyTx.Status,
		Type:      legacyTx.Type,
		CreatedAt: legacyTx.CreatedAt,
		ServiceID: legacyTx.ServiceID,
	}

		return &models.ServiceProlongResponse{
		Transaction: tx,
	}
}
