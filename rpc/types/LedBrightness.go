package types

type SetLedBrightnessRequest struct {
	Request
	Params SetLedBrightnessParams `json:"params"`
}

type SetLedBrightnessParams struct {
	Brightness float64 `json:"brightness"`
}

type SetLedBrightnessResponse struct {
	Response
	Result []SetLedBrightnessParams `json:"result"`
}
