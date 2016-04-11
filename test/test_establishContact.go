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
		fmt.Println(network.GetOwnID())


	case false:
		fmt.Println("This is a slave elevator")
		//start slave elevator routines
	
		ClientEstablishContact(152)
	}



	fmt.Println("I MADE IT ALL THE WAY HERE")


}


func MasterEstablishContact(slaveAddress1, slaveAddress2 int, send_channel, receive_channel chan network.Packet) int {

	err := network.Network_Init(network.LocalListenPort, network.BroadcastListenPort,
        1024, send_channel, receive_channel)
	if err != nil {
		//return err
		fmt.Println("Something went wrong!")
	}

	masterMsg := network.Packet{Receiver_address: "broadcast", Sender_address: string(network.GetOwnID()),
				    Data: []byte("Testmsg"), Length:7}


	numberOfSlavesFound := 0
	firstSlaveFound := ""

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Looking for slaves\n")

		if numberOfSlavesFound == 2 {break}
		send_channel <- masterMsg

		rcvMsg := <- receive_channel
		fmt.Println("Received message from "+rcvMsg.Sender_address)


		if rcvMsg.Sender_address != firstSlaveFound && rcvMsg.Sender_address != string(network.GetOwnID())  {
			firstSlaveFound = rcvMsg.Sender_address
			numberOfSlavesFound++
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


	var masterIP string
	

	for {
		//time.Sleep(1 * time.Second)
		fmt.Println("Looking for Master...")
		rcvMsg := <- receive_channel
		//if rcvMsg.Sender_address == masterIP {break}
		masterIP = rcvMsg.Sender_address
		break

	}

	fmt.Println("Master discovered!")


	clientMsg := network.Packet{Receiver_address: masterIP+":"+strconv.Itoa(network.LocalListenPort),
				    Sender_address: string(network.GetOwnID()), Data: []byte("Testmsg"), Length:7}

	for i:=1; i<10; i++ {
		time.Sleep(1* time.Second)
		send_channel <- clientMsg
	}


}
