package main

import (
	"../driver"
    "../config"
    //"../queue"
    "../statemachine"
    //"../communication"
    "time"
    //"fmt"
)


//var network_sending_channel = make(chan network.Packet)
//var network_receive_channel = make(chan network.Packet)



var network_sending_channel = make(chan config.InfoPackage)

var order_button_pressed_channel = make(chan config.OrderButton)
var floor_reached_channel = make(chan int)

var client_timed_out = make(chan network.ID)

func main() {



	driver.Elev_Init()


	go driver.Order_Button_Poller(order_button_pressed_channel)
	go driver.Floor_Poller(floor_reached_channel)

	//go PollingHandler()
	//go NetworkHandler()



	// ch := statemachine.Internal_channels{
	// 	New_order:     make(chan bool),
	// 	At_floor: make(chan int),
	// 	Direction:     make(chan int),
	// 	Floor_lamp:    make(chan int),
	// 	Door_open_lamp:     make(chan bool),
	// 	Package_out:  network_sending_channel,
	// }
	// //fsm.Init(ch, floor)
	// ch.At_floor <- 5


	for {
		time.Sleep(1 * time.Second)		
	}

}








func listenForSlaveTimeout(ID network.ID, timer *time.Timer, timeOutChan chan network.ID) {
	for {
		select {
		case <- timer.C:
			timeOutChan <- ip
		}

	}
}












