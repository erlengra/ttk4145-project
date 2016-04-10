package main

import (
	"../driver"
	"fmt"
	"time"
	"../config"
	"strconv"
)






func main() {


	tmp := driver.Elev_Init()

	switch tmp {
	case 1:
		fmt.Println("Elevator initiated succesfully")
	case 0:
		fmt.Println("Error during elevator initiation!")
	}



	channel_button_pressed := make(chan config.OrderButton)

	go driver.Order_Button_Poller(channel_button_pressed)


	for {

		time.Sleep(1 * time.Second)
		select{
		case tmpButton := <- channel_button_pressed:
		
			fmt.Println("Button of type "+strconv.Itoa(int(tmpButton.Type))+" pressed at floor "+strconv.Itoa(tmpButton.Floor)+"\n")
		//default:
		//	fmt.Println("Nothing....")


		}

	}	




}
