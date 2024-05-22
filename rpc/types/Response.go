package types

type Response struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

func (r Response) GetID() string {
	return r.ID
}
