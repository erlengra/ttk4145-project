// Package fsm implements a finite state machine for the behaviour of a lift.
// The lift runs based on a queue stored and managed by the queue package.
//
// There are three states:
// Idle: Lift is stationary, at a floor, door closed, awaiting orders.
// Moving: Lift is moving, can be between floors or at a floor going past it.
// Door open: Lift is at a floor with the door open.
//
// And three events:
// New order: A new order is added to the queue.
// Floor reached: The lift reaches a floor.
// Door timeout: The door timer times out (the door should close).
package statemachine

import (
	def "../config"
	"log"
	"../queue"
	"time"
	"../network"
	"fmt"
)

const (
	idle int = iota
	moving
	doorOpen
)

var state int
var floor int
var dir def.MotorDirection

type Channels struct {
	// Events
	NewOrder     chan bool
	FloorReached chan int
	doorTimeout  chan bool
	// Hardware interaction
	MotorDir  chan def.MotorDirection
	FloorLamp chan int
	DoorLamp  chan bool
	// Door timer
	doorTimerReset chan bool
	// Network interaction
	OutgoingMsg chan network.Packet
}

func Init(ch Channels, startFloor int) {
	state = idle
	dir = def.STOP_Direction
	floor = startFloor

	ch.doorTimeout = make(chan bool)
	ch.doorTimerReset = make(chan bool)

	go doorTimer(ch.doorTimeout, ch.doorTimerReset)
	go run(ch)

	//go orderReceiver(ch)

	log.Println(def.ColG, "FSM initialised.", def.ColN)
}

func orderReceiver(ch Channels) {
	for{
		select {
		case <- ch.NewOrder:
			fmt.Println("Tysk")
			eventNewOrder(ch)
		
		}
	}
}

func run(ch Channels) {
	for {
		select {
		case <-ch.NewOrder:
			fmt.Println("Hmmmm")
			eventNewOrder(ch)
		case floor := <-ch.FloorReached:
			eventFloorReached(ch, floor)
		case <-ch.doorTimeout:
			eventDoorTimeout(ch)
		}
	}
}


// func run(ch Channels) {
// 	for {
// 		//fmt.Println("This runnihng")
// 		select {
// 		//case <-ch.NewOrder:
// 		//	fmt.Println("Tysk")
// 		//	eventNewOrder(ch)
// 		case floor := <-ch.FloorReached:
// 			fmt.Println("fRANSK")
// 			eventFloorReached(ch, floor)
// 		//case <-ch.doorTimeout:
// 		//	eventDoorTimeout(ch)
// 		}
// 	}
// }

func eventNewOrder(ch Channels) {
	log.Printf("%sEVENT: New order in state %v.%s", def.ColY, stateString(state), def.ColN)
	switch state {
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
			ch.DoorLamp <- true
			state = doorOpen
		} else {
			ch.MotorDir <- dir
			state = moving
		}
	case moving:
		// Ignore.
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
		}
	default:
		def.CloseConnectionChan <- true
		//def.Restart.Run()
		log.Fatalf(def.ColR, "This state doesn't exist", def.ColN)
	}
}

func eventFloorReached(ch Channels, newFloor int) {
	log.Printf("%sEVENT: Floor %d reached in state %s.%s", def.ColY, newFloor+1, stateString(state), def.ColN)
	floor = newFloor
	ch.FloorLamp <- floor
	switch state {
	case moving:
		fmt.Println("MOTOR STOPPED!")

		if queue.ShouldStop(floor, dir) {
			ch.doorTimerReset <- true
			queue.RemoveOrdersAt(floor, ch.OutgoingMsg)
			ch.DoorLamp <- true
			dir = def.STOP_Direction
			ch.MotorDir <- dir
			state = doorOpen
		}
	default:
		def.CloseConnectionChan <- true
		//def.Restart.Run()
		log.Fatalf("%sMakes no sense to arrive at a floor in state %s%s.\n", def.ColR, stateString(state), def.ColN)
	}
}

func eventDoorTimeout(ch Channels) {
	log.Printf("%sEVENT: Door timeout in state %s.%s", def.ColY, stateString(state), def.ColN)
	switch state {
	case doorOpen:
		ch.DoorLamp <- false
		dir = queue.ChooseDirection(floor, dir)
		ch.MotorDir <- dir
		if dir == def.STOP_Direction {
			state = idle
		} else {
			state = moving
		}
	default:
		def.CloseConnectionChan <- true
		//def.Restart.Run()
		log.Fatalln(def.ColR, "Makes no sense to time out when not in state door open", def.ColN)
	}
}




// Direction returns the current direction as stored in the state machine.
func Direction() def.MotorDirection {
	return dir
}

// Floor returns the current floor as stored in the state machine.
// If the lift is between floors, it returns the most recent floor.
func Floor() int {
	return floor
}

func stateString(state int) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case doorOpen:
		return "door open"
	default:
		return "error: bad state"
	}
}

func doorTimer(timeout chan<- bool, reset <-chan bool) {
	const doorOpenTime = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}