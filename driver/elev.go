package driver

import (
	"fmt"
	"../config"
	"time"
)







var lamp_channel_matrix = [config.NUM_FLOORS][config.NUM_BUTTONS] int {
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [config.NUM_FLOORS][config.NUM_BUTTONS] int {
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}







//Returns 0 on failure, 1 on success
func Elev_Init() int{

	//Beware of the logic here, the if and else may have to be switched for proper response!
	if !Io_Init(){
		fmt.Println("The elevator was initialized correctly!")
	} else {
		fmt.Println("The elevator initialization failed!")
		return 0;
	}

	for f := 0; f < config.NUM_FLOORS; f++ {
		var b config.ButtonType
		for b = 0; b < config.NUM_BUTTONS; b++ {
			Elev_Set_Button_Lamp(b,f,0)
		}
	}

	Elev_Set_Stop_Lamp(0)
	Elev_Set_Door_Open_Lamp(0)
	Elev_Set_Floor_Indicator(0)

	//Only works when we assume that the elevator does not start up at the bottom
	ElevDown()
	for Io_Read_Bit(SENSOR_FLOOR1) != true { 
		fmt.Println("On my way to floor 1...")
		time.Sleep(1 * time.Second)
	}
	ElevStop()
	fmt.Println("Elevator initation complete. Ready to start from floor 1!")
	time.Sleep(1 * time.Second)

	return 1;
}

func Elev_Set_Motor_Direction(dirn config.MotorDirection) {

	if dirn == 0 {
		Io_Write_Analog(MOTOR, 0)
	} else if dirn > 0{
		Io_Clear_Bit(MOTORDIR)
		Io_Write_Analog(MOTOR, config.MOTOR_SPEED)
	} else if dirn < 0 {
		Io_Set_Bit(MOTORDIR)
		Io_Write_Analog(MOTOR, config.MOTOR_SPEED)
	}
}

func ElevUp() {
	Elev_Set_Motor_Direction(config.UP_Direction)
}

func ElevStop() {
	Elev_Set_Motor_Direction(config.STOP_Direction)
}

func ElevDown() {
	Elev_Set_Motor_Direction(config.DOWN_Direction)
}


func Elev_Set_Button_Lamp(button config.ButtonType, floor int, value int) {

	//The following should be checked
	//	assert(floor >= 0);
    //	assert(floor < NUM_FLOORS);
    //	assert(button >= 0);
    //	assert(button < NUM_BUTTONS);
	
	if value == 1 {
		Io_Set_Bit(lamp_channel_matrix[floor][int(button)])
	} else {
		Io_Clear_Bit(lamp_channel_matrix[floor][int(button)])
	}

}


func Elev_Set_Floor_Indicator(floor int) {

	//The following should be checked
	//	assert(floor >= 0);
    //	assert(floor < NUM_FLOORS);


	switch floor {
	case 0:
			Io_Clear_Bit(LIGHT_FLOOR_IND1)
			Io_Clear_Bit(LIGHT_FLOOR_IND2)
	case 1:
			Io_Clear_Bit(LIGHT_FLOOR_IND1)
			Io_Set_Bit(LIGHT_FLOOR_IND2)
	case 2:
			Io_Set_Bit(LIGHT_FLOOR_IND1)
			Io_Clear_Bit(LIGHT_FLOOR_IND2)
	case 3:
			Io_Set_Bit(LIGHT_FLOOR_IND1)
			Io_Set_Bit(LIGHT_FLOOR_IND2)
	}
}


func Elev_Set_Door_Open_Lamp(value int) {
	if value == 1{
		Io_Set_Bit(LIGHT_DOOR_OPEN)
	} else {
		Io_Clear_Bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_Set_Stop_Lamp(value int) {
	if value == 1 {
		Io_Set_Bit(LIGHT_STOP)
	} else {
		Io_Clear_Bit(LIGHT_STOP)
	}
}


func Elev_Get_Button_Signal(button config.ButtonType, floor int) int {
	//The following should be checked
	//    assert(floor >= 0);
    //	  assert(floor < NUM_FLOORS);
    //    assert(button >= 0);
    //    assert(button < NUM_BUTTONS);

	if Io_Read_Bit(button_channel_matrix[floor][int(button)]) {
		return 1
	} else {
		return 0
	}
}

func Elev_Get_Floor_Sensor_Signal() int {

	if Io_Read_Bit(SENSOR_FLOOR1) {
		return 0
	} else if Io_Read_Bit(SENSOR_FLOOR2) {
		return 1
	} else if Io_Read_Bit(SENSOR_FLOOR3) {
		return 2
	} else if Io_Read_Bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func Elev_Get_Stop_Signal() bool {
	return Io_Read_Bit(STOP)
}

func Elev_Get_Obstruction_Signal() bool {
	return Io_Read_Bit(OBSTRUCTION)
}


func Order_Button_Poller(polling_chan_button chan config.OrderButton) {
	var buttonType config.ButtonType
	var lastFloorPassed[config.NUM_BUTTONS][config.NUM_FLOORS]int

	for {
		time.Sleep(100 * time.Millisecond)
		for buttonType = config.BUTTON_CALL_UP; buttonType <= config.BUTTON_CALL_COMMAND; buttonType++ {
			for floor := 0; floor < config.NUM_FLOORS; floor++ {
				buttonValue := Elev_Get_Button_Signal(buttonType, floor)
				if buttonValue != 0 && buttonValue != lastFloorPassed[buttonType][floor] {
					polling_chan_button <- config.OrderButton{Type: buttonType, Floor: floor}
				}
				lastFloorPassed[buttonType][floor] = buttonValue
			}
		}
	}
}



func Floor_Poller(polling_chan_floor chan int) {

	var currentFloor, previousFloor int
	previousFloor = -1

	for {
		time.Sleep(100 * time.Millisecond)
		currentFloor = Elev_Get_Floor_Sensor_Signal()
		if currentFloor != previousFloor && currentFloor != -1 {
			previousFloor = currentFloor
			polling_chan_floor <- currentFloor
		}

	}
}





































