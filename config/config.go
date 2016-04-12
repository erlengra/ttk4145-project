package config

import (
	"fmt"
	"os"
)



const NUM_FLOORS = 4
const NUM_ELEVATORS = 3
const NUM_BUTTONS = 3

const MOTOR_SPEED = 2800


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


type OrderButton struct {
    Floor int
    Type  ButtonType
}


func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}




type InfoPackage struct {
	direction bool
	lastPassedFloor int
	targetFloor int
}