package types

type GetPlugNChargeRequest struct {
	Request
}

type GetPlugNChargeResponseParams struct {
	Applied string `json:"applied"`
}

type GetPlugNChargeResponse struct {
	Response
	Result GetPlugNChargeResponseParams `json:"result"`
}
