package main

import (
	"fmt"
	"log"
	"net"
	"os"

	handlers "github.com/pandaAritra/sqliteWireProtocol/Handlers"
)

func main() {

	// find port
	if len(os.Args) < 2 {
		log.Fatalf("port needed")

	}
	port := fmt.Sprintf(":%s", os.Args[1])

	//establish connection
	listener, err := net.Listen("tcp", port) //binds the port
	if err != nil {
		log.Fatalf(" \ncould not bind port ---------------\n%s ", err)
	}
	defer listener.Close()
	fmt.Println("listning on", listener.Addr())
	for {
		// start accepting clients and fire go rutines for eatch client
		// listner is a pipe , client has a struct in the os karnel
		// with ip, port and a specific buf for each client the
		//accept is a blocking func it blocks until a conn is triggred
		//then it creates the socket and goes back to waiting.
		// thats why we bind it in a infinite loop and we do now error
		// out of it expect graceful shutdown

		client, err := listener.Accept()
		if err != nil {
			fmt.Println("waiting")
			continue
		}
		fmt.Println("client is ", client.RemoteAddr()) // now there is a client object that has propaties

		go handlers.LengthPayload(client) //makes eatch client seperate
	}

}
