package communication
//encoding/decoding

import (
    //"../network"
    "encoding/json"
    "fmt"
)


type clientMessage struct {
	direction bool
	lastPassedFloor int
	targetFloor int
}




func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}


func EncodeClientMessage(message clientMessage) []byte {	
	dataStream, err := json.Marshal(message)
	CheckError(err)
	return dataStream
}

func DecodeClientMessage(dataStream []byte) clientMessage {	
     var message clientMessage
     err := json.Unmarshal(dataStream, &message)
     CheckError(err)
     return message
}
