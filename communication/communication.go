package communication
//encoding/decoding

import (
    "../network"
    "../driver"
    "encoding/json"
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

func DecodeClientMessage(packet network.Packet) clientMessage {	
     var message ClientData
     err := json.Unmarshal(b, &message)
     CheckError(err)
     return result
}
