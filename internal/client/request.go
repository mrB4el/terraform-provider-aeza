// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/client/request.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"io"
	"net/http"
	"net/url"
)

type Request struct {
	client      *Client
	method      string
	path        string
	body        interface{}
	queryParams map[string]string
}

func (c *Client) NewRequest(method, path string, body interface{}) *Request {
	return &Request{
		client:      c,
		method:      method,
		path:        path,
		body:        body,
		queryParams: make(map[string]string),
	}
}

func (r *Request) Do(ctx context.Context, result interface{}) error {
	start := time.Now()

	var bodyReader io.Reader
	var requestBody []byte
	var err error

	if r.body != nil {
		requestBody, err = json.Marshal(r.body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(requestBody)
	}

	fullURL := r.client.host + r.path
	if len(r.queryParams) > 0 {
		query := url.Values{}
		for key, value := range r.queryParams {
			query.Add(key, value)
		}
		fullURL = fullURL + "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, r.method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", r.client.apiKey)
	if r.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	r.client.logRequest(req.Method, fullURL, requestBody, req.Header)

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		duration := time.Since(start)
		r.client.logError(req.Method, fullURL, err, duration)
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		duration := time.Since(start)
		r.client.logError(req.Method, fullURL, err, duration)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	duration := time.Since(start)

	r.client.logResponse(req.Method, fullURL, resp.Status, responseBody, resp.Header, duration)

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode}
	}

	if result != nil {
		if err := json.Unmarshal(responseBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
