package http

// Server - Interface implements fast http server for admin service.
type Server interface {
	// Starts fast HTTP server.
	Start() error
	// Stops fast HTTP server.
	Stop() error
}
