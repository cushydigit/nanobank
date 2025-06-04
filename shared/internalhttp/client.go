package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/cushydigit/nanobank/shared/config"
)

var INTERNAL_AUTH_TOKEN string = os.Getenv("INTERNAL_AUTH_TOKEN")

var Client = &http.Client{
	Timeout: config.INTERNAL_TIMEOUT_CLIENT,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.INTERNAL_CONNECTION_TIMEOUT,
			KeepAlive: config.INTERNAL_KEEPALIVE_TIMEOUT,
		}).DialContext,
		MaxIdleConns:        config.INTERNAL_IDLCONNECTION_MAXCOUNT,
		IdleConnTimeout:     config.INTERNAL_IDLCONNECTION_TIMEOUT,
		TLSHandshakeTimeout: config.INTERNAL_TLSHANDSHAKE_TIMEOUT,
	},
}

func DoJSON(ctx context.Context, method, url string, requestBody any, responseBody any) error {
	var body io.Reader
	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// inject internal auth if set
	if INTERNAL_AUTH_TOKEN != "" {
		req.Header.Set("Authorization", "Bearer "+INTERNAL_AUTH_TOKEN)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(resp.Body)
		return &InternalHTTPError{Status: resp.StatusCode, Body: string(errBody)}
	}

	if responseBody != nil {
		return json.NewDecoder(resp.Body).Decode(responseBody)
	}
	return nil
}

type InternalHTTPError struct {
	Status int
	Body   string
}

func (e *InternalHTTPError) Error() string {
	return http.StatusText(e.Status) + ": " + e.Body
}
