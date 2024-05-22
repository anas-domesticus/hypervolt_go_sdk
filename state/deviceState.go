package state

import "github.com/anas-domesticus/hypervolt_go_sdk/types"

type HypervoltDeviceState struct {
	ChargerID                       string
	IsCharging                      bool
	SessionID                       string
	SessionWatthours                float64
	SessionCurrencySpent            float64
	SessionCarbonSavedGrams         float64
	MaxCurrentMilliamps             float64
	CurrentSessionCurrentMilliamps  float64
	CurrentSessionCtCurrent         float64
	CurrentSessionCtPower           float64
	CurrentSessionVoltage           float64
	EVPower                         float64
	HousePower                      float64
	GridPower                       float64
	GenerationPower                 float64
	LEDBrightness                   float64
	LockState                       *types.HypervoltLockState
	ChargeMode                      *types.HypervoltChargeMode
	ReleaseState                    *types.HypervoltReleaseState
	ActivationMode                  *types.HypervoltActivationMode
	ScheduleIntervals               []*types.HypervoltScheduleInterval
	ScheduleTz                      string
	ScheduleType                    string
	CarPlugged                      bool
	ScheduleIntervalsToApply        []*types.HypervoltScheduleInterval
	SessionWatthoursTotalIncreasing float64
	CurrentSessionPower             float64
}
