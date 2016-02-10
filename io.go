package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: ${SRCDIR}/simelev.a /usr/lib/x86_64-linux-gnu/libphobos2.a -lpthread -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"

func io_init(elevatorType int) int {
	return int(C.io_init(C.ElevatorType(elevatorType)))
}

func io_set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func io_read_bit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func io_read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

func io_write_analog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}
