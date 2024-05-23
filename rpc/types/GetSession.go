package types

type GetSession struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Charging         *bool `json:"charging"`
		TrueMilliAmps    *int  `json:"true_milli_amps"`
		Voltage          *int  `json:"voltage"`
		WattHours        *int  `json:"watt_hours"`
		CarbonSavedGrams *int  `json:"carbon_saved_grams"`
		CtCurrent        *int  `json:"ct_current"`
		EvPower          *int  `json:"ev_power"`
		GridPower        *int  `json:"grid_power"`
		HousePower       *int  `json:"house_power"`
		GenerationPower  *int  `json:"generation_power"`
	} `json:"params"`
}
