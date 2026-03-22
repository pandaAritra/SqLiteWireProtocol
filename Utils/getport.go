package utils

import (
	"fmt"
	"log"
	"os"
)

func GetPort() string {
	if len(os.Args) < 2 {
		log.Fatalf("port no not given")
	}
	return fmt.Sprintf(":%s", os.Args[1])
}
