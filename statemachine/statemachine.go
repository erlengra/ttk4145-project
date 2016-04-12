package statemachine

import (
		//def "config"
		//"queue"
		"../config"
	)

type Internal_channels struct{
	New_order		chan bool
	Package_out 		chan config.InfoPackage
	At_floor 		chan int
	Close_door		chan bool
	Floor_lamp		chan int
	Door_open_lamp		chan bool
	Direction		chan int
	Reset_timer		chan bool
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

func Initialize(channel internal_channels, initial_floor int){
	state = idle
	floor = initial_floor
	direction = def.direction_when_stop
	channel.close_door make(chan bool)
	channel.reset_timer make(chan bool)
	go timer(channel.close_door,channel.reset_timer)
	go execute(channel)
}

func Execute(channel internal_channels){
	for {
		select{
		case <- channel.new_order:
			NewOrder(channel)
		case floor := <- channel.at_floor:
			FloorReached(channel, floor)
		case <- channel.close_door:
			CloseDoor(channel)
		}
	}
}

func NewOrder(channel internal_channels){
	switch state {
	case idle:
		direction = queue.selectDirection(floor,channel.package_out)
		if queue.stop_elevator(floor, direction){
			state = stop
			channel.door_open_lamp <- true
			channel.reset_timer <- true
			queue.removeOrder(floor,channel.package_out)
		}
	case moving_down:
	case moving_up:
	case stop:
		if queue.stop_elevator(floor, direction){
			channel.reset_timer <- true
			queue.removeOrder(floor,channel.package_out)
		}
	default:
		println("invalid state detected")
		def.execute()
	}
}

func FloorReached(channel internal_channels){
	floor = new_floor
	channel.floor_lamp <- floor
	switch state {
	case moving_up:
		if queue.stop_elevator(floor, direction) {
			channel.reset_timer <- true
			queue.RemoveOrder(floor, channel.package_out)
			channel.DoorLamp <- true
			direction = def.direction_when_stop
			channel.elev_direction<- direction
			state = stop
		}
	case moving_down:
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
			ch.DoorLamp <- true
			dir = def.DirStop
			ch.MotorDir <- dir
			state = doorOpen
		}
	default:
		println("invalid state detected")
		def.execute()
	}
}

func CloseDoor(channel internal_channels){
	switch state{
	case idle:
	case stop:
		channel.door_open_lamp <- false
		direction = queue.selectDirection(floor, direction)
		channel.elev_direction <- direction
		if direction == def.direction_when_stop{state = idle}
		else{state = moving_down}
	default:
		println("invalid state detected")
		def.execute()
	}
}
