package statemachine

import (
		"../config"
		"../queue"
		//"../communication"
		"fmt"
		"time"
	)

type internal_channels struct{
	new_order		chan bool
	package_out 	chan config.infoPackage
	at_floor 		chan int
	close_door		chan bool
	floor_lamp		chan int
	door_open_lamp	chan bool
	direction		chan int
	reset_timer		chan bool
}



func timer(timeout chan <- bool, reset_timer <- chan bool){
	timer := time.NewTimer(0)
	timer.stop()

	for{
		if <- reset{
			timer.Reset(timer)
		}
		else if <-timer.C{
			timer.stop()
			timeout <- true
		}
	}
}



const (
	idle int = iota
	moving
	stop
)

var state int
var floor int
var elev_direction int

func initialize(channel internal_channels, initial_floor int){
	state = idle
	floor = initial_floor
	direction = def.direction_when_stop
	channel.close_door = make(chan bool)
	channel.reset_timer = make(chan bool)
	go timer(channel.close_door,channel.reset_timer)
	go mainExecute(channel)
}

func mainExecute(channel internal_channels){
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
		direction = queue.SelectDirection(floor,channel.package_out)
		if queue.StopElevator(floor, direction){
			state = stop
			channel.door_open_lamp <- true
			channel.reset_timer <- true
			queue.RemoveOrder(floor,channel.package_out)
		}
	case moving:
	case stop:
		if queue.StopElevator(floor, direction){
			channel.reset_timer <- true
			queue.RemoveOrder(floor,channel.package_out)
		}
	default:
		println("invalid state detected")
		def.run.execute()
	}
}

func FloorReached(channel internal_channels){
	floor = new_floor
	channel.floor_lamp <- floor
	switch state {
	case moving:
		if queue.StopElevator(floor, direction) {
			channel.reset_timer <- true
			queue.RemoveOrder(floor, channel.package_out)
			channel.DoorLamp <- true
			direction = def.direction_when_stop
			channel.elev_direction <- direction
			state = stop
		}
	default:
		println("invalid state detected")
		def.run.execute()
	}
}

func CloseDoor(channel internal_channels){
	switch state{
	case idle:
	case stop:
		channel.door_open_lamp <- false
		direction = queue.SelectDirection(floor, direction)
		channel.elev_direction <- direction
		if direction == def.direction_when_stop{state = idle}
		else{state = moving}
	default:
		println("invalid state detected")
		def.run.execute()
	}
}