package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// SetLocked sends a request to set the locked state of the charger
func (c *Client) SetLocked(locked bool) (*types.SetLockedResponse, error) {
	req := types.SetLockedRequest{
		Request: types.NewRequest("sync.apply"),
		Params: types.SetLockedRequestParams{
			IsLocked: locked,
		},
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.SetLockedResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
