package main
import (
	"fmt"
    "net"
    "os"
)

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

// Small program for listening for the messages from the UDP server
// This program runs a "server", listening on port 30000
// The return message tells us that the ip of the UDP/TCP server is 129.241.187.23

func main() {

	/* Lets prepare a address at any address at port 30000*/   
	ServerAddr, err := net.ResolveUDPAddr("udp",":30000")
	CheckError(err)

	/* Now listen at selected port */
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 1024)
 
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        fmt.Println("Received ",string(buf[0:n]), " from ",addr)
 
        if err != nil {
            fmt.Println("Error: ",err)
        } 
    }

}
