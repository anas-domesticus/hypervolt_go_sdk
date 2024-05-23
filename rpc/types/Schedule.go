package types

import (
	"encoding/json"
	"time"
)

type GetScheduleRequest struct {
	Request
}

type ScheduleSession struct {
	Days        []DayOfWeek `json:"days"`
	EndTime     time.Time   `json:"end_time"`
	Mode        string      `json:"mode"`
	SessionType string      `json:"session_type"`
	StartTime   time.Time   `json:"start_time"`
}

type DayOfWeek int

const (
	MONDAY DayOfWeek = iota
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
	SUNDAY
	ALL
)

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
	ss.Mode = jsonScheduleSession.Mode
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
		switch day {
		case "MONDAY":
			ss.Days = append(ss.Days, MONDAY)
		case "TUESDAY":
			ss.Days = append(ss.Days, TUESDAY)
		case "WEDNESDAY":
			ss.Days = append(ss.Days, WEDNESDAY)
		case "THURSDAY":
			ss.Days = append(ss.Days, THURSDAY)
		case "FRIDAY":
			ss.Days = append(ss.Days, FRIDAY)
		case "SATURDAY":
			ss.Days = append(ss.Days, SATURDAY)
		case "SUNDAY":
			ss.Days = append(ss.Days, SUNDAY)
		case "ALL":
			ss.Days = append(ss.Days, ALL)
		}
	}
	return nil
}

type GetScheduleResponseParams struct { // TODO: make this better, improve types etc
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

type T struct {
	Days        []string `json:"days"`
	EndTime     string   `json:"end_time"`
	Mode        string   `json:"mode"`
	SessionType string   `json:"session_type"`
	StartTime   string   `json:"start_time"`
}
