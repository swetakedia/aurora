package handlers

import (
	"github.com/blocksafe/go/clients/federation"
	"github.com/blocksafe/go/clients/aurora"
	"github.com/blocksafe/go/clients/blocksafetoml"
	"github.com/blocksafe/go/services/bridge/internal/config"
	"github.com/blocksafe/go/services/bridge/internal/db"
	"github.com/blocksafe/go/services/bridge/internal/listener"
	"github.com/blocksafe/go/services/bridge/internal/submitter"
	"github.com/blocksafe/go/support/http"
)

// RequestHandler implements bridge server request handlers
type RequestHandler struct {
	Config               *config.Config                          `inject:""`
	Client               http.SimpleHTTPClientInterface          `inject:""`
	Aurora              aurora.ClientInterface                 `inject:""`
	Database             db.Database                             `inject:""`
	BlocksafeTomlResolver  blocksafetoml.ClientInterface             `inject:""`
	FederationResolver   federation.ClientInterface              `inject:""`
	TransactionSubmitter submitter.TransactionSubmitterInterface `inject:""`
	PaymentListener      *listener.PaymentListener               `inject:""`
}

func (rh *RequestHandler) isAssetAllowed(code string, issuer string) bool {
	for _, asset := range rh.Config.Assets {
		if asset.Code == code && asset.Issuer == issuer {
			return true
		}
	}
	return false
}
