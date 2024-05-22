package types

type SetMaxCurrentRequest struct {
	Request
	Params SetMaxCurrentParams `json:"params"`
}

type SetMaxCurrentParams struct {
	MaxCurrent int `json:"max_current"`
}

type SetMaxCurrentResponse struct {
	Response
	Result []SetMaxCurrentParams `json:"result"`
}
