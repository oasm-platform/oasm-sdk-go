package oasm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Health() (bool, error) {
	resp, err := c.req.Get(c.apiURL + "/api/health")
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
