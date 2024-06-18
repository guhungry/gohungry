package http

import "net/http"

// HTTP method constants
const (
	MethodGet  = http.MethodGet  // GET method for HTTP requests
	MethodPost = http.MethodPost // POST method for HTTP requests
)

// AuthType represents the type of authentication used in HTTP requests
type AuthType string

// Authentication type constants as defined by relevant RFCs
const (
	AuthTypeBasic  AuthType = "Basic"  // Basic authentication as per RFC 7617
	AuthTypeBearer AuthType = "Bearer" // Bearer token authentication as per RFC 6750
)

// Common HTTP header constants
const (
	HeaderAuthorization = "Authorization" // Header for authorization credentials
	HeaderContentType   = "Content-Type"  // Header indicating the media type of the resource
	HeaderAccept        = "Accept"        // Header indicating the media types acceptable for the response
)
