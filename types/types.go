package types

import "time"

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

type HypervoltDeviceState struct {
	ChargerID                       string
	IsCharging                      *bool
	SessionID                       *string
	SessionWatthours                *float64
	SessionCurrencySpent            *float64
	SessionCarbonSavedGrams         *float64
	MaxCurrentMilliamps             *float64
	CurrentSessionCurrentMilliamps  *float64
	CurrentSessionCtCurrent         *float64
	CurrentSessionCtPower           *float64
	CurrentSessionVoltage           *float64
	EVPower                         *float64
	HousePower                      *float64
	GridPower                       *float64
	GenerationPower                 *float64
	LEDBrightness                   *float64
	LockState                       *HypervoltLockState
	ChargeMode                      *HypervoltChargeMode
	ReleaseState                    *HypervoltReleaseState
	ActivationMode                  *HypervoltActivationMode
	ScheduleIntervals               []*HypervoltScheduleInterval
	ScheduleTz                      *string
	ScheduleType                    *string
	CarPlugged                      *bool
	ScheduleIntervalsToApply        []*HypervoltScheduleInterval
	SessionWatthoursTotalIncreasing *float64
	CurrentSessionPower             *float64
}
