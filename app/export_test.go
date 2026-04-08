package app

import "net/http"

// TestHTTPHandler returns the configured HTTP handler after Setup (for integration tests in app_test).
func (e *Connect) TestHTTPHandler() http.Handler {
	if e == nil || e.core == nil {
		return nil
	}
	return e.core
}
