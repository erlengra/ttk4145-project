package queue

import (
		def "configFile"
		"stateMachine"
		"time"
		"fmt"
)

type queue struct {		queue_table [def.floors][numOfButtons]status	}

type status struct {
	active bool
	inactive bool
	addr
	timer
}

var idle = status{active: false, addr: "", timer: nil}
var private queue
var public queue
var backup = make(chan bool)
var new_order chan bool
var update = make(chan bool)


func DisplayQueue() {
	fmt.Printf(def.ColC)
	fmt.Println("private queue 	 |   public queue")
	for f := def.NumOfFloors - 1; f >= 0; f-- {

		private_column := ""
		if local.IsOrder(f, def.BtnUp) {
			private_column += "↑"
		} else {
			private_column += " "
		}
		if local.IsOrder(f, def.BtnInside) {
			private_column += "×"
		} else {
			private_column += " "
		}
		fmt.Printf(private_column)
		if local.IsOrder(f, def.BtnDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}

		public_column := "   "
		if public.IsOrder(f, def.BtnUp) {
			fmt.Printf("↑")
			public_column += "(↑ " + public.matrix[f][def.BtnUp].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		if public.IsOrder(f, def.BtnDown) {
			fmt.Printf("↓")
			public_column += "(↓ " + public.matrix[f][def.BtnDown].addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", public_column)
		fmt.Println()
	}
	fmt.Printf(def.ColN)
}


func Initialize(hold_new_order chan bool, package_out chan def.infoPackage){
	newOrder = hold_new_order
	go update()
	backup(package_out)
}


func AddOrderToPrivate(floor int, button int){
	private.StoreOrder(floor, button, status{true, "", nil})
	new_order <- true
}


func AddOrderToPublic(floor int, button int){
	is_present := PresentInPublicQueue(floor, button)
	public.SetOrder(floor, button, status{true, addr, nil})
	if !is_present {
		go public.startTimer(floor, button)
	}
	update <- true
}


func RemoveOrder(floor int, package_out chan<- def.infoPackage){
	for i := 0; i < def.NumOfButtons; i++ {
		public.StopTimer(floor, i)
		private.SetOrder(floor, i, inactive)
		public.SetOrder(floor, i, inactive)
	}
	package_out <- def.infoPackage{Category: def.order_done, Floor: floor}
}


func StopElevator(floor, direction int) bool {
	return private.StopElevator(floor, direction)
}


func SelectDirection(floor, direction int) int {
	return private.SelectDirection(floor, direction)
}


func PresentInPrivateQueue(floor, button int) bool{
	return private.IsPresent(floor, button)
}


func PresentInPublicQueue(floor, button int) bool{
	return public.IsPresent(floor,button)
}


func ReassignOrders(elevator_dead string, package_out chan<- def.infoPackage) {
	for i := 0; i < def.NumOfFloors; i++ {
		for j := 0; j < def.NumOfButtons; j++ {
			if public.queue_table[i][j].addr == elevator_dead {
				public.SetOrder(i, j, inactive)
				package_out <- def.infoPackage{Category: def.NewOrder, Floor: i, Button: j}
			}
		}
	}
}
