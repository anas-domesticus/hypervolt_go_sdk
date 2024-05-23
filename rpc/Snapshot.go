package rpc

import (
	"encoding/json"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
)

// GetSnapshot sends a request to retrieve the plug and charge state of the charger.
func (c *Client) GetSnapshot() (*types.GetSnapshotResponse, error) {
	req := types.GetSnapshotRequest{
		Request: types.NewRequest("sync.snapshot"),
	}
	rawResp, err := c.sendMessageAndWaitForResponse(req)
	if err != nil {
		return nil, err
	}
	resp := types.GetSnapshotResponse{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
