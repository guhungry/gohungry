package http

import "net/http"

// HTTP method constants
const (
	MethodGet  = http.MethodGet  // GET method for HTTP requests
	MethodPost = http.MethodPost // POST method for HTTP requests
)

// Authentication type constants
const (
	AuthTypeBasic  = "Basic"  // Basic authentication type (RFC 7617)
	AuthTypeBearer = "Bearer" // Bearer token authentication type (RFC 6750)
)

// Common HTTP header constants
const (
	HeaderAuthorization = "Authorization" // Header for authorization credentials
	HeaderContentType   = "Content-Type"  // Header indicating the media type of the resource
	HeaderAccept        = "Accept"        // Header indicating the media types acceptable for the response
)
