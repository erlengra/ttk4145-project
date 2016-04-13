package main

import (
	"../network"
	"../config"
	"../driver"
	"../statemachine"
	"../queue"
	"fmt"
	//"time"
	"strconv"
)

var network_sending_channel = make(chan network.Packet)
var network_receive_channel = make(chan network.Packet)

var order_button_pressed_channel = make(chan config.OrderButton)
var floor_reached_channel = make(chan int)


func main() {
	driver.Elev_Init()
	var init_floor = 0 //first floor?

	go driver.Order_Button_Poller(order_button_pressed_channel)
	go driver.Floor_Poller(floor_reached_channel)


	ch := statemachine.Channels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan config.MotorDirection, 10),
		FloorLamp:    make(chan int, 10),
		DoorLamp:     make(chan bool, 10),
		OutgoingMsg:  network_sending_channel,
	}


	statemachine.Init(ch, init_floor)
	go eventHandler(ch)

	queue.Init(ch.NewOrder, network_sending_channel)


	balboa := make(chan bool)
	<-balboa
}


func eventHandler(ch statemachine.Channels) {
	//buttonChan := pollButtons()
	//floorChan := pollFloors()

	for {
		select {
		case key := <-order_button_pressed_channel:
			switch key.Type {
			case config.BUTTON_CALL_COMMAND:
				queue.AddLocalOrder(key.Floor, key.Type)
				//fmt.Println("Order added locally!")

			case config.BUTTON_CALL_DOWN, config.BUTTON_CALL_UP:
				//handle
				network_sending_channel <- network.Packet{Receiver_address: "129.241.187.155:"+strconv.Itoa(network.LocalListenPort), Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}
				//network_sending_channel <- config.Message{Category: config.NewOrder, Floor: key.Floor, Button: key.Button}
			}
		case floor := <-floor_reached_channel:
			fmt.Println("SOmethinG")
			ch.FloorReached <- floor
		//case message := <-network_receive_channel:
		//	handleMessage(message)
		//case connection := <-deadChan:
		//	handleDeadLift(connection.Addr)
		//case order := <-queue.OrderTimeoutChan:
		//	log.Println(config.ColR, "Order timeout, I will do it myself!", config.ColN)
		//	queue.RemoveRemoteOrdersAt(order.Floor)
		//	queue.AddRemoteOrder(order.Floor, order.Button, config.Laddr)
		case dir := <-ch.MotorDir:
			fmt.Println("A!")
			driver.Elev_Set_Motor_Direction(dir)
		//case floor := <-ch.FloorLamp:
		//	fmt.Println("B!")
		//	driver.Elev_Set_Floor_Indicator(floor)
		case value := <-ch.DoorLamp:
			fmt.Println("C!")
			driver.Elev_Set_Door_Open_Lamp2(value)

		// case <- ch.NewOrder:
		// 	fmt.Println("Data sent into NewOrder")


		}

	}
}


