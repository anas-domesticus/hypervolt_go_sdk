package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// SetLedBrightness sets the brightness of the charger's LEDs.
// The value parameter should be between 0 and 1.
func (c *Client) SetLedBrightness(value float64) (*types.SetLedBrightnessResponse, error) {
	req := types.SetLedBrightnessRequest{
		Request: types.NewRequest("sync.apply"),
		Params: types.SetLedBrightnessParams{
			Brightness: value,
		},
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.SetLedBrightnessResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
