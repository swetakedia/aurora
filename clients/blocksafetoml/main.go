package blocksafetoml

import "net/http"

// BlocksafeTomlMaxSize is the maximum size of blocksafe.toml file
const BlocksafeTomlMaxSize = 5 * 1024

// WellKnownPath represents the url path at which the blocksafe.toml file should
// exist to conform to the federation protocol.
const WellKnownPath = "/.well-known/blocksafe.toml"

// DefaultClient is a default client using the default parameters
var DefaultClient = &Client{HTTP: http.DefaultClient}

// Client represents a client that is capable of resolving a Blocksafe.toml file
// using the internet.
type Client struct {
	// HTTP is the http client used when resolving a Blocksafe.toml file
	HTTP HTTP

	// UseHTTP forces the client to resolve against servers using plain HTTP.
	// Useful for debugging.
	UseHTTP bool
}

type ClientInterface interface {
	GetBlocksafeToml(domain string) (*Response, error)
	GetBlocksafeTomlByAddress(addy string) (*Response, error)
}

// HTTP represents the http client that a stellertoml resolver uses to make http
// requests.
type HTTP interface {
	Get(url string) (*http.Response, error)
}

// Response represents the results of successfully resolving a blocksafe.toml file
type Response struct {
	AuthServer       string `toml:"AUTH_SERVER"`
	FederationServer string `toml:"FEDERATION_SERVER"`
	EncryptionKey    string `toml:"ENCRYPTION_KEY"`
	SigningKey       string `toml:"SIGNING_KEY"`
}

// GetBlocksafeToml returns blocksafe.toml file for a given domain
func GetBlocksafeToml(domain string) (*Response, error) {
	return DefaultClient.GetBlocksafeToml(domain)
}

// GetBlocksafeTomlByAddress returns blocksafe.toml file of a domain fetched from a
// given address
func GetBlocksafeTomlByAddress(addy string) (*Response, error) {
	return DefaultClient.GetBlocksafeTomlByAddress(addy)
}
