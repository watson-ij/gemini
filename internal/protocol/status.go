package protocol

import "fmt"

// StatusCode represents a Gemini protocol status code
type StatusCode int

// Status code categories and specific codes
const (
	// 1x - INPUT
	StatusInput          StatusCode = 10
	StatusSensitiveInput StatusCode = 11

	// 2x - SUCCESS
	StatusSuccess StatusCode = 20

	// 3x - REDIRECT
	StatusRedirectTemporary StatusCode = 30
	StatusRedirectPermanent StatusCode = 31

	// 4x - TEMPORARY FAILURE
	StatusTemporaryFailure StatusCode = 40
	StatusServerUnavailable StatusCode = 41
	StatusCGIError         StatusCode = 42
	StatusProxyError       StatusCode = 43
	StatusSlowDown         StatusCode = 44

	// 5x - PERMANENT FAILURE
	StatusPermanentFailure     StatusCode = 50
	StatusNotFound             StatusCode = 51
	StatusGone                 StatusCode = 52
	StatusProxyRequestRefused  StatusCode = 53
	StatusBadRequest           StatusCode = 59

	// 6x - CLIENT CERTIFICATE REQUIRED
	StatusClientCertificateRequired StatusCode = 60
	StatusCertificateNotAuthorised  StatusCode = 61
	StatusCertificateNotValid       StatusCode = 62
)

// StatusCategory represents the category of a status code
type StatusCategory int

const (
	CategoryInput           StatusCategory = 1
	CategorySuccess         StatusCategory = 2
	CategoryRedirect        StatusCategory = 3
	CategoryTemporaryFailure StatusCategory = 4
	CategoryPermanentFailure StatusCategory = 5
	CategoryClientCertificate StatusCategory = 6
)

// Category returns the category of a status code
func (s StatusCode) Category() StatusCategory {
	return StatusCategory(s / 10)
}

// IsInput returns true if the status code is in the INPUT category (1x)
func (s StatusCode) IsInput() bool {
	return s.Category() == CategoryInput
}

// IsSuccess returns true if the status code is in the SUCCESS category (2x)
func (s StatusCode) IsSuccess() bool {
	return s.Category() == CategorySuccess
}

// IsRedirect returns true if the status code is in the REDIRECT category (3x)
func (s StatusCode) IsRedirect() bool {
	return s.Category() == CategoryRedirect
}

// IsTemporaryFailure returns true if the status code is in the TEMPORARY FAILURE category (4x)
func (s StatusCode) IsTemporaryFailure() bool {
	return s.Category() == CategoryTemporaryFailure
}

// IsPermanentFailure returns true if the status code is in the PERMANENT FAILURE category (5x)
func (s StatusCode) IsPermanentFailure() bool {
	return s.Category() == CategoryPermanentFailure
}

// IsClientCertificate returns true if the status code is in the CLIENT CERTIFICATE category (6x)
func (s StatusCode) IsClientCertificate() bool {
	return s.Category() == CategoryClientCertificate
}

// IsError returns true if the status code indicates an error (4x or 5x)
func (s StatusCode) IsError() bool {
	return s.IsTemporaryFailure() || s.IsPermanentFailure()
}

// String returns a human-readable description of the status code
func (s StatusCode) String() string {
	switch s {
	case StatusInput:
		return "Input Required"
	case StatusSensitiveInput:
		return "Sensitive Input Required"
	case StatusSuccess:
		return "Success"
	case StatusRedirectTemporary:
		return "Temporary Redirect"
	case StatusRedirectPermanent:
		return "Permanent Redirect"
	case StatusTemporaryFailure:
		return "Temporary Failure"
	case StatusServerUnavailable:
		return "Server Unavailable"
	case StatusCGIError:
		return "CGI Error"
	case StatusProxyError:
		return "Proxy Error"
	case StatusSlowDown:
		return "Slow Down"
	case StatusPermanentFailure:
		return "Permanent Failure"
	case StatusNotFound:
		return "Not Found"
	case StatusGone:
		return "Gone"
	case StatusProxyRequestRefused:
		return "Proxy Request Refused"
	case StatusBadRequest:
		return "Bad Request"
	case StatusClientCertificateRequired:
		return "Client Certificate Required"
	case StatusCertificateNotAuthorised:
		return "Certificate Not Authorised"
	case StatusCertificateNotValid:
		return "Certificate Not Valid"
	default:
		// For undefined status codes, use the category default description
		category := s.Category()
		switch category {
		case CategoryInput:
			return fmt.Sprintf("Input Required (%d)", s)
		case CategorySuccess:
			return fmt.Sprintf("Success (%d)", s)
		case CategoryRedirect:
			return fmt.Sprintf("Redirect (%d)", s)
		case CategoryTemporaryFailure:
			return fmt.Sprintf("Temporary Failure (%d)", s)
		case CategoryPermanentFailure:
			return fmt.Sprintf("Permanent Failure (%d)", s)
		case CategoryClientCertificate:
			return fmt.Sprintf("Client Certificate Issue (%d)", s)
		default:
			return fmt.Sprintf("Unknown Status (%d)", s)
		}
	}
}
