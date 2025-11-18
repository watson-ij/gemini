package protocol

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultPort is the default Gemini port
	DefaultPort = "1965"

	// DefaultTimeout is the default connection timeout
	DefaultTimeout = 30 * time.Second

	// MaxRedirects is the maximum number of redirects to follow
	MaxRedirects = 5
)

// Client is a Gemini protocol client
type Client struct {
	// Timeout is the connection timeout
	Timeout time.Duration

	// TLSConfig is the TLS configuration
	// If nil, a default config will be used
	TLSConfig *tls.Config

	// FollowRedirects controls whether redirects are automatically followed
	FollowRedirects bool

	// MaxRedirects is the maximum number of redirects to follow
	MaxRedirects int

	// TOFU is the Trust On First Use certificate verifier
	TOFU *TOFUVerifier
}

// NewClient creates a new Gemini client with default settings
func NewClient() *Client {
	return &Client{
		Timeout:         DefaultTimeout,
		FollowRedirects: true,
		MaxRedirects:    MaxRedirects,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			// InsecureSkipVerify is set to true because we do TOFU verification
			// We'll verify certificates manually in the TOFU verifier
			InsecureSkipVerify: true,
		},
	}
}

// Get performs a GET request to the specified URL
func (c *Client) Get(rawURL string) (*Response, error) {
	return c.get(rawURL, 0)
}

// get is the internal implementation that tracks redirect count
func (c *Client) get(rawURL string, redirectCount int) (*Response, error) {
	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Validate URL scheme
	if u.Scheme != "gemini" {
		return nil, fmt.Errorf("unsupported URL scheme: %s (expected gemini)", u.Scheme)
	}

	// Get host and port
	host := u.Host
	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, DefaultPort)
	}

	// Create TLS connection
	dialer := &net.Dialer{
		Timeout: c.Timeout,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", host, c.TLSConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer func() {
		// We'll only close the connection if there's an error
		// For success responses, the caller is responsible for closing
		if err != nil {
			conn.Close()
		}
	}()

	// Verify certificate with TOFU if available
	if c.TOFU != nil {
		hostname, _, _ := net.SplitHostPort(host)
		if hostname == "" {
			hostname = host
		}

		if err := c.TOFU.VerifyCertificate(hostname, conn.ConnectionState()); err != nil {
			return nil, fmt.Errorf("certificate verification failed: %w", err)
		}
	}

	// Send request
	request := rawURL + "\r\n"
	if _, err := io.WriteString(conn, request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	resp, err := ReadResponse(conn, rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Store TLS state
	// resp.TLSState = &conn.ConnectionState()

	// Handle redirects
	if resp.Status.IsRedirect() {
		// Close the connection since we won't be reading the body
		conn.Close()

		// Check if we should follow redirects
		if !c.FollowRedirects {
			return resp, nil
		}

		// Check redirect limit
		if redirectCount >= c.MaxRedirects {
			return nil, fmt.Errorf("too many redirects (max %d)", c.MaxRedirects)
		}

		// Parse redirect URL
		redirectURL := resp.Meta
		if redirectURL == "" {
			return nil, fmt.Errorf("redirect without URL")
		}

		// Resolve relative URLs
		redirectURL, err = resolveURL(rawURL, redirectURL)
		if err != nil {
			return nil, fmt.Errorf("invalid redirect URL: %w", err)
		}

		// Follow redirect
		return c.get(redirectURL, redirectCount+1)
	}

	// For non-success, non-redirect responses, close the connection
	if !resp.Status.IsSuccess() {
		conn.Close()
	}

	return resp, nil
}

// resolveURL resolves a potentially relative URL against a base URL
func resolveURL(base, ref string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(refURL).String(), nil
}

// Request performs a request with the given URL and options
// This is a more flexible alternative to Get
func (c *Client) Request(rawURL string, opts ...RequestOption) (*Response, error) {
	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c.Get(rawURL)
}

// RequestOption is a function that configures a request
type RequestOption func(*Client)

// WithTimeout sets the connection timeout
func WithTimeout(timeout time.Duration) RequestOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// WithTLSConfig sets the TLS configuration
func WithTLSConfig(config *tls.Config) RequestOption {
	return func(c *Client) {
		c.TLSConfig = config
	}
}

// WithoutRedirects disables automatic redirect following
func WithoutRedirects() RequestOption {
	return func(c *Client) {
		c.FollowRedirects = false
	}
}
