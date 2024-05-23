package types

import (
	"encoding/json"
	"time"
)

type ScheduleSession struct {
	Days        []DayOfWeek         `json:"days"`
	EndTime     time.Time           `json:"end_time"`
	Mode        HypervoltChargeMode `json:"mode"`
	SessionType string              `json:"session_type"`
	StartTime   time.Time           `json:"start_time"`
}

func (ss *ScheduleSession) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var jsonScheduleSession struct {
		Days        []string `json:"days"`
		EndTime     string   `json:"end_time"`
		Mode        string   `json:"mode"`
		SessionType string   `json:"session_type"`
		StartTime   string   `json:"start_time"`
	}
	if err := json.Unmarshal(data, &jsonScheduleSession); err != nil {
		return err
	}
	ss.Mode = HypervoltChargeMode(jsonScheduleSession.Mode)
	ss.SessionType = jsonScheduleSession.SessionType

	t, err := time.Parse("15:04", jsonScheduleSession.StartTime)
	if err != nil {
		return err
	}
	ss.StartTime = t
	t, err = time.Parse("15:04", jsonScheduleSession.EndTime)
	if err != nil {
		return err
	}
	ss.EndTime = t

	ss.Days = []DayOfWeek{}
	for _, day := range jsonScheduleSession.Days {
		ss.Days = append(ss.Days, DayOfWeek(day))
	}
	return nil
}

func (ss *ScheduleSession) MarshalJSON() ([]byte, error) {
	var jsonScheduleSession struct {
		Days        []string `json:"days"`
		EndTime     string   `json:"end_time"`
		Mode        string   `json:"mode"`
		SessionType string   `json:"session_type"`
		StartTime   string   `json:"start_time"`
	}
	jsonScheduleSession.Mode = string(ss.Mode)
	jsonScheduleSession.SessionType = ss.SessionType
	jsonScheduleSession.StartTime = ss.StartTime.Format("15:04")
	jsonScheduleSession.EndTime = ss.EndTime.Format("15:04")
	for _, day := range ss.Days {
		jsonScheduleSession.Days = append(jsonScheduleSession.Days, string(day))
	}

	return json.Marshal(jsonScheduleSession)
}

type GetScheduleRequest struct {
	Request
}

type GetScheduleResponseParams struct {
	Applied struct {
		Enabled   bool              `json:"enabled"`
		IsDefault bool              `json:"is_default"`
		Sessions  []ScheduleSession `json:"sessions"`
		Type      string            `json:"type"`
	} `json:"applied"`
	Pending struct {
	} `json:"pending"`
}

type GetScheduleResponse struct {
	Response
	Result GetScheduleResponseParams `json:"result"`
}

type SetScheduleRequest struct {
	Request
	Params SetScheduleRequestParams `json:"params"`
}

type SetScheduleRequestParams struct {
	Enabled   bool              `json:"enabled"`
	IsDefault bool              `json:"is_default"`
	Sessions  []ScheduleSession `json:"sessions"`
	Type      string            `json:"type"`
}

type SetScheduleResponse struct {
	Response
	Result GetScheduleResponseParams `json:"result"`
}

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
		Enabled   bool              `json:"enabled"`
		IsDefault bool              `json:"is_default"`
		Sessions  []ScheduleSession `json:"sessions"`
		Type      string            `json:"type"`
	} `json:"applied"`
	Pending struct {
	} `json:"pending"`
}

// {"id":"1716470495334144","jsonrpc":"2.0","result":{"applied":{"enabled":true,"is_default":false,"sessions":[{"days":["THURSDAY"],"end_time":"05:00","mode":"boost","session_type":"recurring","start_time":"04:00"}],"type":"hypervolt"},"pending":{}}}

type T2 struct {
	Type      string `json:"type"`
	Tz        string `json:"tz"`
	Intervals [][]struct {
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
		Seconds int `json:"seconds"`
	} `json:"intervals"`
}
