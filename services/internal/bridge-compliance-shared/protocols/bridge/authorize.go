package bridge

// AuthorizeRequest represents request made to /authorize endpoint of bridge server
type AuthorizeRequest struct {
	AccountID string `form:"account_id" valid:"required,blocksafe_accountid"`
	AssetCode string `form:"asset_code" valid:"required,blocksafe_asset_code"`
}

func (r AuthorizeRequest) Validate(params ...interface{}) error {
	// No custom validations
	return nil
}
