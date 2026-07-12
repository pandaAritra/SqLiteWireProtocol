package utils

import (
	"fmt"
	"log"
	"os"
)

// GetPort parses the port from command line arguments.
func GetPort() string {
	if len(os.Args) < 2 {
		log.Fatalf("error: port number not specified")
	}
	return fmt.Sprintf(":%s", os.Args[1])
}
