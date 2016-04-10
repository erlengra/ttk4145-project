package main

import (
    "./network"
    //"./driver"
    //"./communication"
    //"os"
    "flag"
    "fmt"
    //"./master"
    "strconv"
    "time"
)

//testing
func main() {

	//To indicate that a process should be the master process, run "go run main.go -master=true"
	//If it should be a slave elevator just run "go run main.go,", as it defaults to false"
	isMasterElevator := flag.Bool("master", false, "Set to true to indicate master elevator")
	flag.Parse()

	network_send_channel := make(chan network.Packet)
	network_receive_channel := make(chan network.Packet)
	


	switch *isMasterElevator {
	case true:
		
		fmt.Println("This is the master elevator")
		//Start master elevator routines

		//temp := master.MasterEstablishContact(152, 125)
		MasterEstablishContact(144,142, network_send_channel, network_receive_channel)
		//fmt.Println(network.GetOwnID())


	case false:
		fmt.Println("This is a slave elevator")
		//start slave elevator routines
	
		ClientEstablishContact(152)
	}



	fmt.Println("I MADE IT ALL THE WAY HERE")


}


//Modify this so that it gets communication channels as parameters/args.
//This means that the channels can be used in other functions, after contact has been established.
//Declaring them in main.go is probably the easiest
func MasterEstablishContact(slaveAddress1, slaveAddress2 int, send_channel, receive_channel chan network.Packet) int {

	//send_channel := make(chan network.Packet)
	//receive_channel := make(chan network.Packet)

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
	//masterIP := "129.241.187."+strconv.Itoa(masterAddress)

	var masterIP string
	

	//Waiting for message from Master before continuing
	for {
		//time.Sleep(1 * time.Second)
		fmt.Println("Looking for Master...")
		rcvMsg := <- receive_channel
		//if rcvMsg.Sender_address == masterIP {break}
		masterIP = rcvMsg.Sender_address
		break

	}
	//At this point we know that the Master elevator is up

	fmt.Println("Master discovered!")


	clientMsg := network.Packet{Receiver_address: masterIP+":"+strconv.Itoa(network.LocalListenPort),
				    Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}


	//Sending 10 messages just to be sure, this needs to be changed
	for i:=1; i<10; i++ {
		time.Sleep(1* time.Second)
		send_channel <- clientMsg
	}


}
