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
