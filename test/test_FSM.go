package main

import (
	"../driver"
    "../config"
    "../queue"
    "../statemachine"
    "../communication"
)


//var network_sending_channel = make(chan network.Packet)
//var network_receive_channel = make(chan network.Packet)



var network_sending_channel = make(chan config.InfoPackage)

var order_button_pressed_channel = make(chan config.OrderButton)
var floor_reached_channel = make(chan int)

func main() {



	driver.Elev_Init()


	go driver.Order_Button_Poller(order_button_pressed_channel)
	go driver.Floor_Poller(floor_reached_channel)

	//go PollingHandler()
	//go NetworkHandler()



	ch := statemachine.internal_channels{
		new_order:     make(chan bool),
		at_floor: make(chan int),
		direction:     make(chan int),
		floor_lamp:    make(chan int),
		door_open_lamp:     make(chan bool),
		package_out:  network_sending_channel,
	}
	//fsm.Init(ch, floor)




}





















