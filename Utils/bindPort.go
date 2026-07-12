package utils

import (
	"log"
	"net"
)

// BindPort binds to the specified TCP port.
func BindPort(port string) net.Listener {
	listener, err := net.Listen("tcp", port) // binds the port
	if err != nil {
		log.Fatalf(" \ncould not bind port ---------------\n%s ", err)
	}

	return listener
}
