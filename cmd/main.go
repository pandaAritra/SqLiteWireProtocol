package main

import (
	"fmt"
	"log"

	handlers "github.com/pandaAritra/sqliteWireProtocol/Handlers"
	utils "github.com/pandaAritra/sqliteWireProtocol/Utils"
)

func main() {

	port := utils.GetPort()
	listener := utils.BindPort(port)
	log.Println("listning on", listener.Addr())

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("waiting")
			continue
		}
		fmt.Println("client is ", client.RemoteAddr()) // now there is a client object that has propaties

		go handlers.LengthPayload(client) //makes eatch client seperate
	}

}
