package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// SetScheduleEnabled sets the brightness of the charger's LEDs.
// The value parameter should be between 0 and 1.
func (c *Client) SetScheduleEnabled(value bool) (*types.SetScheduleEnabledResponse, error) {
	req := types.SetScheduleEnabledRequest{
		Request: types.NewRequest("schedule.set"),
		Params: types.SetScheduleEnabledRequestParams{
			Enabled: value,
		},
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.SetScheduleEnabledResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
