package state_classes


// StateClass refers to the type of sensor state being tracked for statistics
//
// https://developers.home-assistant.io/docs/core/entity/sensor/#available-state-classes
type StateClass string

const (
	Measurement = "measurement"
    Total="total"
    TotalIncreasing="total_increasing"
	None        = ""
)
