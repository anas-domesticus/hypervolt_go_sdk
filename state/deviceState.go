package state

import (
	"time"
)

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
	LockState                       HypervoltLockState
	ChargeMode                      HypervoltChargeMode
	ReleaseState                    HypervoltReleaseState
	ActivationMode                  HypervoltActivationMode
	ScheduleIntervals               []HypervoltScheduleInterval
	ScheduleTz                      string
	ScheduleType                    string
	CarPlugged                      bool
	ScheduleIntervalsToApply        []HypervoltScheduleInterval
	SessionWatthoursTotalIncreasing float64
	CurrentSessionPower             float64
	RandomStart                     bool
}

type HypervoltLockState int

const (
	UNLOCKED HypervoltLockState = iota
	PENDING_LOCK
	LOCKED
)

type HypervoltChargeMode int

const (
	BOOST HypervoltChargeMode = iota
	ECO
	SUPER_ECO
)

type HypervoltActivationMode int

const (
	PLUG_AND_CHARGE HypervoltActivationMode = iota
	SCHEDULE
)

type HypervoltReleaseState int

const (
	DEFAULT HypervoltReleaseState = iota
	RELEASED
)

type HypervoltDayOfWeek int

const (
	MONDAY HypervoltDayOfWeek = 1 << iota
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
	SUNDAY
	ALL = MONDAY | TUESDAY | WEDNESDAY | THURSDAY | FRIDAY | SATURDAY | SUNDAY
)

const NUM_SCHEDULE_INTERVALS = 4

type HypervoltScheduleInterval struct {
	StartTime  time.Time
	EndTime    time.Time
	ChargeMode HypervoltChargeMode
	DaysOfWeek HypervoltDayOfWeek
}
