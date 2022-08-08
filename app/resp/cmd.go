package resp

import (
	"log"
)

func Ping() []byte {
	return encodeSimpleString("PONG")
}

func Echo(v Value) []byte {
	data, err := encode(&v)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		return []byte{}
	}
	return data
}
