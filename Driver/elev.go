package driver

import "fmt"

const N_FLOORS = 4
const N_BUTTONS = 3

const MOTOR_SPEED = 2800


type elev_button_type_t int
const (
	BUTTON_CALL_UP = iota // 0
	BUTTON_CALL_DOWN      // 1
	BUTTON_COMMAND        // 2
)

type elev_motor_direction_t int
const (
    UP_Direction_t =  1
    DOWN_Direction = -1
    STOP_Direction =  0
)



var lamp_channel_matrix = [N_FLOORS][N_BUTTONS] int {
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS] int {
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}







//Returns 0 on failure, 1 on success
func Elev_Init() int{

	//Beware of the logic here, the if and else may have to be switched!
	if !Io_Init(){
		fmt.Println("The elevator was initialized correctly!")
	} else {
		fmt.Println("The elevator initialization failed!")
		return 0;
	}

	for f := 0; f < N_FLOORS; f++ {
		var b elev_button_type_t
		for b = 0; b < N_BUTTONS; b++ {
			Elev_Set_Button_Lamp(b,f,0)
		}
	}

	Elev_Set_Stop_Lamp(0)
	Elev_Set_Door_Open_Lamp(0)
	Elev_Set_Floor_Indicator(0)

	return 1;
}

func Elev_Set_Motor_Direction(dirn elev_motor_direction_t) {

	if dirn == 0 {
		Io_Write_Analog(MOTOR, 0)
	} else if dirn > 0{
		Io_Clear_Bit(MOTORDIR)
		Io_Write_Analog(MOTOR, MOTOR_SPEED)
	} else if dirn < 0 {
		Io_Set_Bit(MOTORDIR)
		Io_Write_Analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_Set_Button_Lamp(button elev_button_type_t, floor int, value int) {

	//The following should be checked
	//	assert(floor >= 0);
    //	assert(floor < N_FLOORS);
    //	assert(button >= 0);
    //	assert(button < N_BUTTONS);
	
	if value == 1 {
		Io_Set_Bit(lamp_channel_matrix[floor][int(button)])
	} else {
		Io_Clear_Bit(lamp_channel_matrix[floor][int(button)])
	}

}


func Elev_Set_Floor_Indicator(floor int) {

	//The following should be checked
	//	assert(floor >= 0);
    //	assert(floor < N_FLOORS);


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


func Elev_Get_Button_Signal(button elev_button_type_t, floor int) int {
	//The following should be checked
	//    assert(floor >= 0);
    //	  assert(floor < N_FLOORS);
    //    assert(button >= 0);
    //    assert(button < N_BUTTONS);

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