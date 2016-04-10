package config



const NUM_FLOORS = 4
const NUM_ELEVATORS = 3
const NUM_BUTTONS = 3



type ButtonType int
const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_CALL_COMMAND
)

type MotorDirection int
const (
	UP_Direction = 1
	DOWN_Direction = -1
	STOP_Direction = 0
)
