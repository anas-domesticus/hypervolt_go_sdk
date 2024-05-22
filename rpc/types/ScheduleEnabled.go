package types

type SetScheduleEnabledRequest struct {
	Request
	Params SetScheduleEnabledRequestParams `json:"params"`
}

type SetScheduleEnabledRequestParams struct {
	Enabled bool `json:"enabled"`
}

type SetScheduleEnabledResponse struct {
	Response
	Result SetScheduleEnabledResponseParams `json:"result"`
}

type SetScheduleEnabledResponseParams struct {
	Applied struct {
		Enabled   bool `json:"enabled"`
		IsDefault bool `json:"is_default"`
		Sessions  []struct {
			Days        []string `json:"days"`
			EndTime     string   `json:"end_time"`
			Mode        string   `json:"mode"`
			SessionType string   `json:"session_type"`
			StartTime   string   `json:"start_time"`
		} `json:"sessions"`
		Type string `json:"type"`
	} `json:"applied"`
	Pending struct {
	} `json:"pending"`
}
