package utils

import (
	"log"
	"net"
)

func BindPort(port string) net.Listener {
	listener, err := net.Listen("tcp", port) //binds the port
	if err != nil {
		defer listener.Close()
		log.Fatalf(" \ncould not bind port ---------------\n%s ", err)
	}

	return listener

}
