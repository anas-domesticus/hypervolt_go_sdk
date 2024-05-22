package types

func NewLoginRequest(token string) *LoginRequest {
	return &LoginRequest{
		Request: NewRequest("login"),
		Params: LoginRequestParams{
			Token:   token,
			Version: 2,
		},
	}
}

type LoginRequest struct {
	Request
	Params LoginRequestParams `json:"params"`
}

type LoginRequestParams struct {
	Token   string `json:"token"`
	Version int    `json:"version"`
}

type LoginResponse struct {
	Response
	Result struct {
		Authenticated bool `json:"authenticated"`
	} `json:"result"`
}
