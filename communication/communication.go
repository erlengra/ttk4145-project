package communication
//encoding/decoding

import (
    //"../network"
    "encoding/json"
    "fmt"
    "../config"
)






func EncodeClientMessage(message config.InfoPackage) []byte {	
	dataStream, err := json.Marshal(message)
	CheckError(err)
	return dataStream
}

func DecodeClientMessage(dataStream []byte) config.InfoPackage {	
     var message config.InfoPackage
     err := json.Unmarshal(dataStream, &message)
     CheckError(err)
     return message
}
