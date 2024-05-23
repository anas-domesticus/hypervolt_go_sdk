package types

type HypervoltChargeMode string

const (
	BOOST     HypervoltChargeMode = "boost"
	ECO       HypervoltChargeMode = "eco"
	SUPER_ECO HypervoltChargeMode = "super_eco"
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
