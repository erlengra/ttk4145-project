package main

import (
	"../network"
	"fmt"
	"time"
	"log"
	"strconv"
)



func main() {

	var localIP = network.GetOwnID()

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort, 1024, send_channel, receive_channel)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)

		msg := network.Packet{Receiver_address: "129.241.187.142:"+strconv.Itoa(network.LocalListenPort), Sender_address: string(localIP),
				      Data: []byte("Testmsg"), Length:7}



		fmt.Println("Sending------")
		send_channel <- msg
		network.PrintPacket(msg)
	}
}
