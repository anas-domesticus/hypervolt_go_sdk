package types

type HypervoltChargeMode string

const (
	BOOST     HypervoltChargeMode = "boost"
	ECO       HypervoltChargeMode = "eco"
	SUPER_ECO HypervoltChargeMode = "super_eco"
)

type HypervoltSessionType string

const (
	RECURRING HypervoltSessionType = "recurring"
)

type DayOfWeek string

const (
	MONDAY    DayOfWeek = "MONDAY"
	TUESDAY   DayOfWeek = "TUESDAY"
	WEDNESDAY DayOfWeek = "WEDNESDAY"
	THURSDAY  DayOfWeek = "THURSDAY"
	FRIDAY    DayOfWeek = "FRIDAY"
	SATURDAY  DayOfWeek = "SATURDAY"
	SUNDAY    DayOfWeek = "SUNDAY"
	ALL       DayOfWeek = "ALL"
)
