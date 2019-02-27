package hal

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkBuilder(t *testing.T) {
	// Link expansion test
	check := func(href string, base string, expectedResult string) {
		lb := LinkBuilder{mustParseURL(base)}
		result := lb.expandLink(href)
		assert.Equal(t, expectedResult, result)
	}

	check("/root", "", "/root")
	check("/root", "//blocksafe.org", "//blocksafe.org/root")
	check("/root", "https://blocksafe.org", "https://blocksafe.org/root")
	check("//else.org/root", "", "//else.org/root")
	check("//else.org/root", "//blocksafe.org", "//else.org/root")
	check("//else.org/root", "https://blocksafe.org", "//else.org/root")
	check("https://else.org/root", "", "https://else.org/root")
	check("https://else.org/root", "//blocksafe.org", "https://else.org/root")
	check("https://else.org/root", "https://blocksafe.org", "https://else.org/root")

	// Regression: ensure that parameters are not escaped
	check("/accounts/{id}", "https://blocksafe.org", "https://blocksafe.org/accounts/{id}")
}

func mustParseURL(base string) *url.URL {
	if base == "" {
		return nil
	}

	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	return u
}
