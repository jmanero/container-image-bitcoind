package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Endpoint for the target bitcoind REST API
var Endpoint = url.URL{Scheme: "http", Host: "127.0.0.1:8332", Path: "/rest"}

// GetBlockchainInfo retrieves the current BlockchainInfo struct from the configured node's REST API
func GetBlockchainInfo(ctx context.Context) (info BlockchainInfo, err error) {
	uri := Endpoint.JoinPath("chaininfo.json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer Drain(res.Body)

	if res.StatusCode != http.StatusOK {
		return info, fmt.Errorf("Unhandled response %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&info)

	return
}

// Drain is a helper that ensures a ReadCloser is fully read and closed
func Drain(rc io.ReadCloser) error {
	_, err := io.Copy(io.Discard, rc)
	rc.Close()
	return err
}
