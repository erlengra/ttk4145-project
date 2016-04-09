package main

import (
    //"./network"
    //"./driver"
    //"./communication"
    //"os"
    "flag"
    "fmt"
)

//testing
func main() {

    //To indicate that a process should be the master process, run "go run main.go -master=true"
    //If it should be a slave elevator just run "go run main.go,", as it defaults to false"
    isMasterElevator := flag.Bool("master", false, "Set to true to indicate master elevator")
    flag.Parse()


    switch *isMasterElevator {
    case true:
        fmt.Println("This is the master elevator")
        //Start master elevator routines

    case false:
        fmt.Println("This is a slave elevator")
        //start slave elevator routines
    }


}
