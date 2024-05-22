package types

import (
	"fmt"
	"time"
)

type Request struct {
	ID     string `json:"id"`
	Method string `json:"method"`
}

func NewRequest(method string) Request {
	return Request{
		ID:     getNextMessageID(),
		Method: method,
	}
}

func (r Request) GetID() string {
	return r.ID
}

func getNextMessageID() string {
	timestamp := time.Now().UTC().UnixNano() / int64(time.Microsecond)
	return fmt.Sprintf("%d", timestamp)
}
