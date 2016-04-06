package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"


func Io_Init() bool {
	return bool(int(C.io_init()) != 1)
}

func Io_Set_Bit(channel int) {
    C.io_set_bit(C.int(channel))
}

func Io_Clear_Bit(channel int) {
    C.io_clear_bit(C.int(channel))
}

func Io_Read_Bit(channel int) bool {
    return bool(int(C.io_read_bit(C.int(channel))) !=0 )
}

func Io_Read_Analog(channel int) int {
    return int(C.io_read_analog(C.int(channel)))
}

func Io_Write_Analog(channel int, value int) {
    C.io_write_analog(C.int(channel),C.int(value))
}