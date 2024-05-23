package state

import "github.com/anas-domesticus/hypervolt_go_sdk/types"

type HypervoltDeviceState struct {
	ChargerID                       string
	IsCharging                      bool
	SessionID                       string
	SessionWatthours                int
	SessionCurrencySpent            float64
	SessionCarbonSavedGrams         int
	MaxCurrentMilliamps             int
	CurrentSessionCurrentMilliamps  int
	CurrentSessionCtCurrent         int
	CurrentSessionCtPower           int
	CurrentSessionVoltage           int
	EVPower                         int
	HousePower                      int
	GridPower                       int
	GenerationPower                 int
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
