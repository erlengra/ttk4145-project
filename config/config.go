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


//////////////////////////////////////////////////////


// Local IP address
var Laddr string

var SyncLightsChan = make(chan bool)
var CloseConnectionChan = make(chan bool)


const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m"    // Grey (neutral)