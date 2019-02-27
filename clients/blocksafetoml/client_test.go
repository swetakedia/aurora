package blocksafetoml

import (
	"net/http"
	"strings"
	"testing"

	"github.com/blocksafe/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientURL(t *testing.T) {
	//HACK:  we're testing an internal method rather than setting up a http client
	//mock.

	c := &Client{UseHTTP: false}
	assert.Equal(t, "https://blocksafe.org/.well-known/blocksafe.toml", c.url("blocksafe.org"))

	c = &Client{UseHTTP: true}
	assert.Equal(t, "http://blocksafe.org/.well-known/blocksafe.toml", c.url("blocksafe.org"))
}

func TestClient(t *testing.T) {
	h := httptest.NewClient()
	c := &Client{HTTP: h}

	// happy path
	h.
		On("GET", "https://blocksafe.org/.well-known/blocksafe.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation"`,
		)
	stoml, err := c.GetBlocksafeToml("blocksafe.org")
	require.NoError(t, err)
	assert.Equal(t, "https://localhost/federation", stoml.FederationServer)

	// blocksafe.toml exceeds limit
	h.
		On("GET", "https://toobig.org/.well-known/blocksafe.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation`+strings.Repeat("0", BlocksafeTomlMaxSize)+`"`,
		)
	stoml, err = c.GetBlocksafeToml("toobig.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "blocksafe.toml response exceeds")
	}

	// not found
	h.
		On("GET", "https://missing.org/.well-known/blocksafe.toml").
		ReturnNotFound()
	stoml, err = c.GetBlocksafeToml("missing.org")
	assert.EqualError(t, err, "http request failed with non-200 status code")

	// invalid toml
	h.
		On("GET", "https://json.org/.well-known/blocksafe.toml").
		ReturnJSON(http.StatusOK, map[string]string{"hello": "world"})
	stoml, err = c.GetBlocksafeToml("json.org")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "toml decode failed")
	}
}
