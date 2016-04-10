package main

import (
	"../network"
	"fmt"
	"time"
	"log"
)



func main() {

	//var localIP = network.GetOwnID()

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort, 1024, send_channel, receive_channel)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Receiving----")

		select{
		case rcvMsg := <- receive_channel:
			network.PrintPacket(rcvMsg)

		}

	}
}
