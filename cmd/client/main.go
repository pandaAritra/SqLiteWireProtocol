package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	utils "github.com/pandaAritra/sqliteWireProtocol/Utils"
	"github.com/pandaAritra/sqliteWireProtocol/models"
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

	readResponse(conn)
}

func readResponse(conn net.Conn) {
	// 1. Read message type
	msgTypeBuf := make([]byte, 1)
	_, err := io.ReadFull(conn, msgTypeBuf)
	if err != nil {
		fmt.Printf("Error: failed to read response message type: %v\n", err)
		return
	}
	msgType := msgTypeBuf[0]

	// 2. Read payload length
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(conn, lenBuf)
	if err != nil {
		fmt.Printf("Error: failed to read response payload length: %v\n", err)
		return
	}
	payloadLength := binary.BigEndian.Uint32(lenBuf)

	// 3. Read payload
	payloadBuf := make([]byte, payloadLength)
	_, err = io.ReadFull(conn, payloadBuf)
	if err != nil {
		fmt.Printf("Error: failed to read response payload: %v\n", err)
		return
	}

	// 4. Handle response based on msgType
	if msgType != 0x03 {
		fmt.Printf("Error: unexpected response message type 0x%02x: %s\n", msgType, string(payloadBuf))
		return
	}

	var response models.Response
	if err := json.Unmarshal(payloadBuf, &response); err != nil {
		fmt.Printf("Error: failed to decode response JSON: %v\n", err)
		return
	}

	if response.Error != "" {
		fmt.Printf("Error: %s\n", response.Error)
		return
	}

	if len(response.Columns) == 0 {
		fmt.Println("Query executed successfully (no rows returned).")
		return
	}

	// Print as table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	
	// Headers
	fmt.Fprintln(w, strings.Join(response.Columns, "\t"))
	
	// Separators
	seps := make([]string, len(response.Columns))
	for i, col := range response.Columns {
		seps[i] = strings.Repeat("-", len(col))
	}
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Rows
	for _, row := range response.Rows {
		rowStr := make([]string, len(row))
		for i, val := range row {
			if val == nil {
				rowStr[i] = "NULL"
			} else {
				rowStr[i] = fmt.Sprintf("%v", val)
			}
		}
		fmt.Fprintln(w, strings.Join(rowStr, "\t"))
	}
	w.Flush()
	fmt.Println()
}
