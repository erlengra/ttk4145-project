package main
 
import (
    "fmt"
    "net"
    "time"
    //"strconv"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
 

// Simple program for sending a message to the server

func main() {


    //Establish socket for sending
    ServerAddr,err := net.ResolveUDPAddr("udp","129.241.187.23:20017")
    CheckError(err) 
    sendingConn, err := net.DialUDP("udp", nil, ServerAddr)
    CheckError(err)
    defer sendingConn.Close()

    //Establish socket for receiving
    // receivingConn, err := net.ListenUDP("udp", ServerAddr)
    // CheckError(err)
    // defer receivingConn.Close()


    // Send a message to the server
    msg := "Test1"
    buf := []byte(msg)
    _,err = sendingConn.Write(buf)

    if err != nil {
        fmt.Println(msg, err)
    } else {
        fmt.Println("Sent message to server")
    }

    //Receive a message from the server
    receive_buf := make([]byte, 1024)

    n,addr,err := receivingConn.ReadFromUDP(receive_buf)
    fmt.Println("Received ",string(receive_buf[0:n]), " from ",addr)
    if err != nil {
        fmt.Println("Error: ",err)
    } 
      
    //Wait
    time.Sleep(time.Second * 1)
    
}
