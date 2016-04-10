package master

import (
	"../network"
	"strconv"
	"time"
	"fmt"
)


//Modify this so that it gets communication channels as parameters/args.
//This means that the channels can be used in other functions, after contact has been established.
//Declaring them in main.go is probably the easiest
func MasterEstablishContact(slaveAddress1, slaveAddress2 int) int {

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort,
        1024, send_channel, receive_channel)
	if err != nil {
		//return err
		fmt.Println("Something went wrong!")
	}


	slave1IP := "129.241.187."+strconv.Itoa(slaveAddress1)
	slave2IP := "129.241.187."+strconv.Itoa(slaveAddress2)

	//numberOfSlavesFound = 0
	slave1Found := false
	slave2Found := false
	masterMsg := network.Packet{Receiver_address: "broadcast", Sender_address: string(network.GetOwnID()),
				    Data: []byte("Testmsg"), Length:7}


	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Looking for slaves\n")

		if slave1Found && slave2Found {break}
		send_channel <- masterMsg

		rcvMsg := <- receive_channel
		if rcvMsg.Sender_address == slave1IP {
			slave1Found = true
			fmt.Println("Slave elevator number 1, at spot "+strconv.Itoa(slaveAddress1)+"found!")
		} else if rcvMsg.Sender_address == slave2IP {
			slave2Found = true
			fmt.Println("Slave elevator number 2, at spot "+strconv.Itoa(slaveAddress2)+"found!")
		}

	}


	fmt.Println("Both elevators seem to be up!")
	return 2

}


//Should probably take in channels as arguments as well.
//Or declare them "globally"...?

//TO-DO: Add functionality to time out if master is not found?

func ClientEstablishContact(masterAddress int) {

	send_channel := make(chan network.Packet)
	receive_channel := make(chan network.Packet)

	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort,
        1024, send_channel, receive_channel)
	if err != nil {
		//return err
		fmt.Println("Something went wrong!")
	}

	//masterDiscovered := false
	masterIP := "129.241.187."+strconv.Itoa(masterAddress)


	//////////////////////////////////////////
	masterIP = "10.22.69.27"
	/////////////////////////////////////////
	

	//Waiting for message from Master before continuing
	for {
		//time.Sleep(1 * time.Second)
		rcvMsg := <- receive_channel
		if rcvMsg.Sender_address == masterIP {break}
	}
	//At this point we know that the Master elevator is up

	fmt.Println("Master discovered!")


	clientMsg := network.Packet{Receiver_address: masterIP+strconv.Itoa(network.LocalListenPort),
				    Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}


	//Sending 10 messages just to be sure, this needs to be changed
	for i:=1; i<10; i++ {
		time.Sleep(1* time.Second)
		send_channel <- clientMsg
	}


}





























