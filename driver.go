package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"

// Let the direction parameter for Set_motor_direction be an int. Create enum or something similar later
//	-1 = DOWN
//	 0 = STOP
//	 1 = UP
//
// Let the buttonType parameter for setting lamps be an int. Create enum or something similar later
//    BUTTON_CALL_UP = 0,
//    BUTTON_CALL_DOWN = 1,
//    BUTTON_COMMAND = 2

//************************************************************
func ioInit() int {
	return int(C.io_init(1))
}

func ioSetBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func ioClearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func ioReadBit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func ioReadAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

func ioWriteAnalog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

//*************************************************************
func Init() {
	C.elev_init()
}

func Set_motor_direction(direction int) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(direction))
}

func Set_button_lamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Get_button_signal(button int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

/*func Get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}

func Get_obstruction() int {
	return int(C.elev_get_obstruction_signal())
}*/
