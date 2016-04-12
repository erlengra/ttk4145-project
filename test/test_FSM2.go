package main

import (
	"../network"
	"../config"
	"../driver"
	"../statemachine"

)

var network_sending_channel = make(chan network.Packet)
var network_receive_channel = make(chan network.Packet)

var order_button_pressed_channel = make(chan config.OrderButton)
var floor_reached_channel = make(chan int)


func main() {
	driver.Elev_Init()

	go driver.Order_Button_Poller(order_button_pressed_channel)
	go driver.Floor_Poller(floor_reached_channel)

	channel := statemachine.Internal_channels {
		New_order: make(chan bool)
		Package_out: make(chan config.InfoPackage)
		At_floor : make(chan int)



	}







}
