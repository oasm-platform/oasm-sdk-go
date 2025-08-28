package oasm

// Option represents a functional option for configuring the Client.
type Option func(*Client)

// WithApiURL sets the base API URL for the Client.
func WithApiURL(apiUrl string) func(*Client) {
	return func(c *Client) {
		c.apiURL = apiUrl
	}
}

// WithApiKey sets the API key for the Client.
func WithApiKey(apiKey string) func(*Client) {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}
