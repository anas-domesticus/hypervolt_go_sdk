package rpc

import (
	"encoding/json"
	"errors"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

func (c *Client) Authenticate(token string) (*types.LoginResponse, error) {
	req := types.NewLoginRequest(token)
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.LoginResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Result.Authenticated {
		return nil, errors.New("authentication failed")
	}
	return &resp, nil
}
