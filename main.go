package main

import (
    "./network"
    "./config"
    "./driver"
    //"./communication"
    //"os"
    "flag"
    "fmt"
    //"./master"
    "strconv"
    "time"
)

var network_sending_channel = make(chan network.Packet)
var network_receive_channel = make(chan network.Packet)

var order_button_pressed_channel = make(chan config.OrderButton)
var floor_reached_channel = make(chan int)

var masterIP = "129.241.187.150"



func main() {

	//To indicate that a process should be the master process, run "go run main.go -master=true"
	//If it should be a slave elevator just run "go run main.go,", as it defaults to false"
	isMasterElevator := flag.Bool("master", false, "Set to true to indicate master elevator")
	flag.Parse()

	switch *isMasterElevator {
	case true:
		fmt.Println("This is the master elevator")
		//Start master elevator routines

	case false:
		fmt.Println("This is a slave elevator")
		//start slave elevator routines
	}



	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort, 1024,
			    network_sending_channel, network_receive_channel)
	config.CheckError(err)


	driver.Elev_Init()


	go driver.Order_Button_Poller(order_button_pressed_channel)
	go driver.Floor_Poller(floor_reached_channel)

	go PollingHandler()
	go NetworkHandler()





	



	for {
		time.Sleep( 1 * time.Second)
	}

}





func PollingHandler() {
	for {
		select {
		case orderButton := <- order_button_pressed_channel:
			//Do...

			fmt.Println("Button of type "+strconv.Itoa(int(orderButton.Type))+" pressed at floor "+
				strconv.Itoa(orderButton.Floor)+"\n")

			network_sending_channel <- network.Packet{Receiver_address: masterIP+":"+strconv.Itoa(network.LocalListenPort),
		        Sender_address: string(network.GetOwnID()), Data: []byte("testMsg"), Length:7}



		case floorReached := <- floor_reached_channel:
			//Do..
			fmt.Println("Reached floor number "+strconv.Itoa(floorReached+1))

		}
	}
}



func NetworkHandler() {
	for {
		message := <- network_receive_channel
		fmt.Println("A button was pressed at "+message.Sender_address+"\n")
	}
} 





















