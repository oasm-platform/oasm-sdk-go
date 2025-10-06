package oasm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Health checks the health status of the API by sending a GET request
// to the "/api/health" endpoint.
//
// It returns true if the response body equals "OK" and the status code is 200 (OK).
// Otherwise, it returns false along with an error if something goes wrong
// (e.g., network error, unexpected status code, or invalid response body).
func (c *Client) Health() (bool, error) {
	resp, err := c.Get(c.apiURL + "/api/health")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(body)) == "OK", nil
}
