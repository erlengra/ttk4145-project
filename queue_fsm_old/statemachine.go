package statemachine

import (
		//def "config"
		//"queue"
		"../config"
		//"time"
	)

type Internal_channels struct{	
	New_order		chan bool
	At_floor 		chan int
	Close_door		chan bool

	Direction		chan config.MotorDirection
	Floor_lamp		chan int
	Door_open_lamp	chan bool

	Reset_timer		chan bool

	Package_out 	chan config.InfoPackage
}

const (
	idle int = iota
	moving_up
	moving_down
	stop
)

var state int
var floor int
var elev_direction int

func Initialize(channel Internal_channels, initial_floor int){
	state = idle
	floor = initial_floor
	elev_direction = config.STOP_Direction

	channel.Close_door = make(chan bool)
	channel.Reset_timer = make(chan bool)
	//go timer(channel.Close_door,channel.Reset_timer)
	//go Execute(channel)
}

// func Execute(channel Internal_channels){
// 	for {
// 		select{
// 		case <- channel.New_order:
// 			NewOrder(channel)
// 		case floor := <- channel.At_floor:
// 			FloorReached(channel, floor)
// 		case <- channel.Close_door:
// 			CloseDoor(channel)
// 		}
// 	}
// }

// func NewOrder(channel Internal_channels){
// 	switch state {
// 	case idle:
// 		elev_direction = queue.selectDirection(floor,channel.package_out)
// 		if queue.stop_elevator(floor, direction){
// 			state = stop
// 			channel.door_open_lamp <- true
// 			channel.reset_timer <- true
// 			queue.removeOrder(floor,channel.package_out)
// 		}
// 	case moving_down:
// 	case moving_up:
// 	case stop:
// 		if queue.stop_elevator(floor, direction){
// 			channel.reset_timer <- true
// 			queue.removeOrder(floor,channel.package_out)
// 		}
// 	default:
// 		println("invalid state detected")
// 		def.execute()
// 	}
// }

// func FloorReached(channel internal_channels){
// 	floor = new_floor
// 	channel.floor_lamp <- floor
// 	switch state {
// 	case moving_up:
// 		if queue.stop_elevator(floor, direction) {
// 			channel.reset_timer <- true
// 			queue.RemoveOrder(floor, channel.package_out)
// 			channel.DoorLamp <- true
// 			direction = def.direction_when_stop
// 			channel.elev_direction<- direction
// 			state = stop
// 		}
// 	case moving_down:
// 		if queue.ShouldStop(floor, dir) {
// 			ch.doorTimerReset <- true
// 			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
// 			ch.DoorLamp <- true
// 			dir = def.DirStop
// 			ch.MotorDir <- dir
// 			state = doorOpen
// 		}
// 	default:
// 		println("invalid state detected")
// 		def.execute()
// 	}
// }

// func CloseDoor(channel internal_channels){
// 	switch state{
// 	case idle:
// 	case stop:
// 		channel.door_open_lamp <- false
// 		direction = queue.selectDirection(floor, direction)
// 		channel.elev_direction <- direction
// 		if direction == def.direction_when_stop{
// 			state = idle
// 		} else{state = moving_down}
// 	default:
// 		println("invalid state detected")
// 		def.execute()
// 	}
// }



// func timer(timeout chan bool, reset_timer chan bool){
// 	timer := time.NewTimer(0)
// 	timer.stop()

// 	for{
// 		if <- reset{
// 			timer.Reset(timer)
// 		}
// 		else if <-timer.C{
// 			timer.stop()
// 			timeout <- true
// 		}
// 	}
// }