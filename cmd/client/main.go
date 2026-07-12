package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	utils "github.com/pandaAritra/sqliteWireProtocol/Utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("SQLite Wire Protocol Client")
		fmt.Println("\nUsage:")
		fmt.Println("  Interactive REPL:  go run cmd/client/main.go <port_or_address>")
		fmt.Println("  Single Query:      go run cmd/client/main.go <port_or_address> \"<sql_query>\"")
		fmt.Println("\nExamples:")
		fmt.Println("  go run cmd/client/main.go 8000")
		fmt.Println("  go run cmd/client/main.go localhost:8000 \"SELECT * FROM users;\"")
		os.Exit(1)
	}

	target := os.Args[1]
	// If only a port number is given (e.g., 8000), default to localhost:8000
	if _, err := strconv.Atoi(target); err == nil {
		target = "localhost:" + target
	}

	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("Error connecting to server at %s: %v\n", target, err)
		os.Exit(1)
	}
	defer conn.Close()

	// Single Query Mode
	if len(os.Args) >= 3 {
		query := strings.Join(os.Args[2:], " ")
		sendQuery(conn, query)
		return
	}

	// Interactive REPL Mode
	fmt.Printf("Connected to SQLite Wire Protocol Server at %s\n", target)
	fmt.Println("Type your SQL queries below. Type 'exit' or 'quit' to close the connection.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("sqlite-wire> ")
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		lowerQuery := strings.ToLower(query)
		if lowerQuery == "exit" || lowerQuery == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		sendQuery(conn, query)
	}
}

func sendQuery(conn net.Conn, query string) {
	// The wire protocol uses type 0x02 for Query
	packetPtr := utils.MakePKT(0x02, query)
	if packetPtr == nil {
		fmt.Println("Error: failed to encode query packet")
		return
	}

	_, err := conn.Write(*packetPtr)
	if err != nil {
		fmt.Printf("Error: failed to send query to server: %v\n", err)
		return
	}
	fmt.Printf("Query sent successfully: %s\n", query)
}
