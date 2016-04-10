package master

import (
	"../network"
	"strconv"



)


//Modify this so that it gets communication channels as parameters/args.
//This means that the channels can be used in other functions, after contact has been established.
func establishContact(slaveAddress1, slaveAddress2 int) int {

	//Init network
	//...

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(network.localListenPort, network.broadcastListenPort, 1024, send_channel, receive_channel)

	slave1IP = "129.241.187."+strconv.Itoa(slaveAddress1)+":"+strconv.Itoa(network.master_port)
	slave2IP = "129.241.187."+strconv.Itoa(slaveAddress2)+":"+strconv.Itoa(network.master_port)




}




