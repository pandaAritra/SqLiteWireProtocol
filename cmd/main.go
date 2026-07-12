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
	log.Println("listening on", listener.Addr())

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("waiting")
			continue
		}
		fmt.Println("client is", client.RemoteAddr()) // now there is a client object that has properties

		go handlers.LengthPayload(client) // makes each client separate
	}
}
