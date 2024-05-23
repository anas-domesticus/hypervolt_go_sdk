package types

import (
	"encoding/json"
	"reflect"
)

type GetSnapshotRequest struct {
	Request
}

type GetSnapshotResponseJSON struct {
	Response
	Result []GetSnapshotParams `json:"result"`
}

type GetSnapshotResponse struct {
	Response
	Result GetSnapshotParams `json:"result"`
}

func (gsr *GetSnapshotResponse) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var jsonGetSnapshotResponse struct {
		Response
		Result []GetSnapshotParams `json:"result"`
	}
	if err := json.Unmarshal(data, &jsonGetSnapshotResponse); err != nil {
		return err
	}
	gsr.Response = jsonGetSnapshotResponse.Response
	valGS := reflect.ValueOf(&gsr.Result).Elem()

	for i := 0; i < valGS.NumField(); i++ {
		valJGS := reflect.ValueOf(&jsonGetSnapshotResponse.Result[i]).Elem()
		fieldGS := valGS.Field(i)
		fieldJGS := valJGS.Field(i)

		if !fieldJGS.IsNil() {
			fieldGS.Set(fieldJGS)
		}
	}
	return nil
}

type GetSnapshotParams struct {
	Brightness   *float64             `json:"brightness,omitempty"`
	LockState    *string              `json:"lock_state,omitempty"`
	ReleaseState *string              `json:"release_state,omitempty"`
	MaxCurrent   *int                 `json:"max_current,omitempty"`
	CtFlags      *int                 `json:"ct_flags,omitempty"`
	SolarMode    *HypervoltChargeMode `json:"solar_mode,omitempty"`
	Features     *[]string            `json:"features,omitempty"`
	RandomStart  *bool                `json:"random_start,omitempty"`
	EffectName   *string              `json:"effect_name,omitempty"`
}
