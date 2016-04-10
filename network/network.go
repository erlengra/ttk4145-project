package network

import (
	"fmt"
	"net"
	"log"
	"strconv"
)

var local_address, broadcast_address *net.UDPAddr


var LocalListenPort int = 10012
var BroadcastListenPort int = 20012


type Packet struct {
	Receiver_address	string
	Sender_address string
	Data	[]byte
	Length	int
}

func Network_Init(localListenPort, broadcastListenPort, msgSize int, sendCh, receiveCh chan Packet) (err error) {
	broadcast_address, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}

	tempConn, err := net.DialUDP("udp4", nil, broadcast_address)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr() //Not this
	local_address, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	local_address.Port = localListenPort

	fmt.Println(local_address)

	localListenConn, err := net.ListenUDP("udp4", local_address)
	if err != nil {
		return err
	}

	broadcastListenConn, err := net.ListenUDP("udp", broadcast_address)
	if err != nil {
		localListenConn.Close()
		return err
	}

	go udpReceiveServer(localListenConn, broadcastListenConn, msgSize, receiveCh)
	go udpTransmitServer(localListenConn, broadcastListenConn, sendCh)

	return err
}

func udpTransmitServer(lconn, bconn *net.UDPConn, sendCh chan Packet) {
	var err error

	for {
		msg := <-sendCh
		if msg.Receiver_address == "broadcast" {
			_, err = lconn.WriteToUDP(msg.Data, broadcast_address)
		} else {
			raddr, err := net.ResolveUDPAddr("udp", msg.Receiver_address)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			_, err = lconn.WriteToUDP(msg.Data, raddr)
		}
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
}

func udpReceiveServer(lconn, bconn *net.UDPConn, msgSize int, receiveCh chan Packet) {
	bconnReceiveCh := make(chan Packet)
	lconnReceiveCh := make(chan Packet)

	go udpConnectionReader(lconn, msgSize, lconnReceiveCh)
	go udpConnectionReader(bconn, msgSize, bconnReceiveCh)

	for {
		select {
		case buf := <-bconnReceiveCh:
			receiveCh <- buf
		case buf := <-lconnReceiveCh:
			receiveCh <- buf
		}
	}
}

func udpConnectionReader(conn *net.UDPConn, msgSize int, receiveCh chan Packet) {
	for {
		buf := make([]byte, msgSize)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		receiveCh <- Packet{Receiver_address: raddr.String(), Sender_address: string(GetOwnID()),
	                            Data: buf[:n], Length: n}
	}
}



type ID string

func GetSenderID(sender *net.UDPAddr) ID {
	return ID(sender.IP.String())
}

func GetOwnID() ID {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ID(ipnet.IP.String())
				}
			}
		}
	}
	return "127.0.0.1"
}


func PrintPacket(msg Packet) {
	fmt.Printf("msg: \n\t raddr = %s \n\t data = %s \n\t length = %v\n", msg.Receiver_address, string(msg.Data), msg.Length)
}




