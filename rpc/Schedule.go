package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// GetSchedule sends a request to retrieve the schedule of the charger.
func (c *Client) GetSchedule() (*types.GetScheduleResponse, error) {
	req := types.GetScheduleRequest{
		Request: types.NewRequest("schedules.get"),
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.GetScheduleResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetSchedule(sessions []types.ScheduleSession, enable bool) (*types.SetScheduleResponse, error) {
	req := types.SetScheduleRequest{
		Request: types.NewRequest("schedule.set"),
		Params: types.SetScheduleRequestParams{
			Enabled:   enable,
			IsDefault: false,
			Type:      "hypervolt",
			Sessions:  sessions,
		},
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.SetScheduleResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetScheduleEnabled sets whether the charger is in Schedule mode
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
