package oasm

import (
	"fmt"

	"github.com/bytedance/sonic"
)

// APIError represents an error response returned by the API.
// It implements the error interface so it can be used as a standard Go error.
type APIError struct {
	Message    string `json:"message"`    // Human-readable error message
	ErrType    string `json:"error"`      // Type or category of the error
	StatusCode int    `json:"statusCode"` // HTTP status code
}

// Error implements the error interface for APIError.
// It provides a formatted string representation of the API error.
func (e *APIError) Error() string {
	return fmt.Sprintf("api error: %s (%s, %d)", e.Message, e.ErrType, e.StatusCode)
}

// ErrorResponse attempts to unmarshal the response body into an APIError.
// If successful, it returns an *APIError as error; otherwise, it returns the unmarshal error.
func ErrorResponse(body []byte) error {
	var apiErr APIError
	if err := sonic.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("failed to parse error response: %w", err)
	}
	return &apiErr
}
