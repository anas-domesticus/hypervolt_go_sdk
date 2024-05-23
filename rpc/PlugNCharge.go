package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// GetPlugNCharge sends a request to retrieve the plug and charge state of the charger.
func (c *Client) GetPlugNCharge() (*types.GetPlugNChargeResponse, error) {
	req := types.GetPlugNChargeRequest{
		Request: types.NewRequest("plugncharge.get"),
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.GetPlugNChargeResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
