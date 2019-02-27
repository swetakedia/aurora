// Package blocksafecore is a client library for communicating with an
// instance of blocksafe-core using through the server's HTTP port.
package blocksafecore

import "net/http"

// SetCursorDone is the success message returned by blocksafe-core when a cursor
// update succeeds.
const SetCursorDone = "Done"

// HTTP represents the http client that a blocksafecore client uses to make http
// requests.
type HTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

// confirm interface conformity
var _ HTTP = http.DefaultClient
