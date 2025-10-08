package redfish

import (
	"time"
)

// RedfishResponse represents a response from a Redfish API call
type RedfishResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Data       interface{}         `json:"data"`
}

// AuthMethod represents Redfish authentication methods
type AuthMethod string

const (
	AuthMethodBasic   AuthMethod = "basic"
	AuthMethodSession AuthMethod = "session"
)

// ClientConfig represents configuration for a Redfish client
type ClientConfig struct {
	Address            string
	Port               int
	Username           string
	Password           string
	AuthMethod         AuthMethod
	TLSServerCACert    string
	InsecureSkipVerify bool
	MaxRetries         int
	InitialDelay       time.Duration
	MaxDelay           time.Duration
	BackoffFactor      float64
	Jitter             bool
}

// DefaultClientConfig returns default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Port:               443,
		AuthMethod:         AuthMethodSession,
		InsecureSkipVerify: false,
		MaxRetries:         3,
		InitialDelay:       time.Second,
		MaxDelay:           60 * time.Second,
		BackoffFactor:      2.0,
		Jitter:             true,
	}
}

// DiscoveredHost represents a host discovered via SSDP
type DiscoveredHost struct {
	Address     string `json:"address"`
	ServiceRoot string `json:"service_root"`
}

// RedfishError represents a Redfish-specific error
type RedfishError struct {
	Message string
	Code    int
}

func (e *RedfishError) Error() string {
	return e.Message
}

// IsRetryable determines if an error is retryable
func IsRetryable(err error) bool {
	if redfishErr, ok := err.(*RedfishError); ok {
		// Don't retry validation errors
		if redfishErr.Code >= 400 && redfishErr.Code < 500 {
			return false
		}
	}

	// Retry network and server errors
	return true
}
