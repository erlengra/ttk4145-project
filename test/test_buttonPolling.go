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

	var tmpButton config.OrderButton

	for {
		tmpButton <- channel_button_pressed
		
		fmt.Println("Button of type "+strconv.Itoa(tmpButton.Type)+" pressed at floor "
			   +strconv.Itoa(tmpButton.Floor)+"\n")


	}	




}
