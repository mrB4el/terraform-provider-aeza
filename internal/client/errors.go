// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package client

import "fmt"

type APIError struct {
	StatusCode int
	Slug       string
	Message    string
	Data       interface{}
}

func (e *APIError) Error() string {
	if e.Message == "" && e.Slug == "" {
		return fmt.Sprintf("API error %d", e.StatusCode)
	}
	return fmt.Sprintf("API Error %d: %s (%s)", e.StatusCode, e.Message, e.Slug)
}

type ErrorResponse struct {
	Error struct {
		Slug    string      `json:"slug"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	} `json:"error"`
}
