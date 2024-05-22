package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// SetMaxCurrent sets the maximum current value in milliamperes (mA).
func (c *Client) SetMaxCurrent(value int) (*types.SetMaxCurrentResponse, error) {
	req := types.SetMaxCurrentRequest{
		Request: types.NewRequest("sync.apply"),
		Params: types.SetMaxCurrentParams{
			MaxCurrent: value,
		},
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.SetMaxCurrentResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
