package main

import (
	"../network"
	"fmt"
	"time"
	"log"
)



func main() {

	var localIP = network.GetOwnID()

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(20001, 20002, 1024, send_channel, receive_channel)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)

		msg := network.Packet{Receiver_address: string(localIP)+":20001", Data: []byte("Testmsg"), Length:7}
		fmt.Println("Sending------")
		send_channel <- msg
		network.PrintPacket(msg)
		fmt.Println("Receiving----")
		rcvMsg := <- receive_channel
		network.PrintPacket(rcvMsg)
	}
}
