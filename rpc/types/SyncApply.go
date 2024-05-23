package types

type SyncApply struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []struct {
		Brightness   *float64  `json:"brightness,omitempty"`
		LockState    *string   `json:"lock_state,omitempty"`
		ReleaseState *string   `json:"release_state,omitempty"`
		MaxCurrent   *int      `json:"max_current,omitempty"`
		CtFlags      *int      `json:"ct_flags,omitempty"`
		SolarMode    *string   `json:"solar_mode,omitempty"`
		Features     *[]string `json:"features,omitempty"`
		RandomStart  *bool     `json:"random_start,omitempty"`
		EffectName   *string   `json:"effect_name,omitempty"`
	} `json:"params"`
}
