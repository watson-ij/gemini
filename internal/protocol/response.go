package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Response represents a Gemini protocol response
type Response struct {
	// Status is the response status code
	Status StatusCode

	// Meta contains the meta string (MIME type for success, error message for failures, etc.)
	Meta string

	// Body contains the response body (nil for non-success responses)
	Body io.ReadCloser

	// TLSState contains information about the TLS connection
	// TLSState *tls.ConnectionState

	// URL is the URL that was requested
	URL string
}

// ParseResponseHeader parses the response header line
// Format: <STATUS><SPACE><META><CR><LF>
func ParseResponseHeader(line string) (StatusCode, string, error) {
	// Remove trailing CR/LF if present
	line = strings.TrimRight(line, "\r\n")

	// Find the first space
	spaceIdx := strings.IndexByte(line, ' ')
	if spaceIdx == -1 {
		// No space means no meta - this is technically invalid but we'll be lenient
		status, err := parseStatus(line)
		if err != nil {
			return 0, "", err
		}
		return status, "", nil
	}

	// Parse status code
	statusStr := line[:spaceIdx]
	status, err := parseStatus(statusStr)
	if err != nil {
		return 0, "", err
	}

	// Everything after the first space is meta
	meta := line[spaceIdx+1:]

	// Meta should be at most 1024 bytes per spec
	if len(meta) > 1024 {
		return 0, "", fmt.Errorf("meta string too long: %d bytes (max 1024)", len(meta))
	}

	return status, meta, nil
}

// parseStatus parses a status code string
func parseStatus(s string) (StatusCode, error) {
	// Status code should be exactly 2 digits
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid status code length: %d (expected 2)", len(s))
	}

	code, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid status code: %s", s)
	}

	// Validate status code range (10-69)
	if code < 10 || code > 69 {
		return 0, fmt.Errorf("status code out of range: %d (expected 10-69)", code)
	}

	return StatusCode(code), nil
}

// ReadResponse reads and parses a Gemini response from a reader
func ReadResponse(r io.Reader, url string) (*Response, error) {
	bufReader := bufio.NewReader(r)

	// Read the header line (terminated by CR LF)
	headerLine, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response header: %w", err)
	}

	// Parse the header
	status, meta, err := ParseResponseHeader(headerLine)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response header: %w", err)
	}

	resp := &Response{
		Status: status,
		Meta:   meta,
		URL:    url,
	}

	// For success responses, the body follows
	// For other responses, there is no body
	if status.IsSuccess() {
		// Don't close the underlying reader - let the caller handle it
		resp.Body = io.NopCloser(bufReader)
	}

	return resp, nil
}

// Close closes the response body if present
func (r *Response) Close() error {
	if r.Body != nil {
		return r.Body.Close()
	}
	return nil
}

// MIMEType returns the MIME type from the meta field for success responses
func (r *Response) MIMEType() string {
	if !r.Status.IsSuccess() {
		return ""
	}

	// The meta field for success responses is: <MIME type>[; <parameters>]
	// We'll just return the MIME type part
	parts := strings.SplitN(r.Meta, ";", 2)
	return strings.TrimSpace(parts[0])
}

// IsGemtext returns true if the response is a gemtext document
func (r *Response) IsGemtext() bool {
	mimeType := r.MIMEType()
	return mimeType == "text/gemini" || mimeType == ""
}

// ReadBody reads the entire response body
func (r *Response) ReadBody() ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("no response body")
	}

	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
