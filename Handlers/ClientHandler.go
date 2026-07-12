package handlers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	mydb "github.com/pandaAritra/sqliteWireProtocol/db"
)

// getDatabasePath resolves the SQLite database path by checking:
// 1. SQLITE_DB_PATH environment variable
// 2. Local relative db/test.db
// 3. Absolute path to current workspace database
// 4. Absolute path to previous database location
func getDatabasePath() string {
	if envPath := os.Getenv("SQLITE_DB_PATH"); envPath != "" {
		return envPath
	}

	paths := []string{
		"db/test.db",
		"/home/panda/Projects/SqLiteWireProtocol/db/test.db",
		"/home/panda/Documents/code/test/SqLiteWireProtocol/db/test.db",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "db/test.db" // Fallback to relative path
}

func EchoClient(client net.Conn) {
	defer client.Close()
	buf := make([]byte, 1024)
	for {
		n, err := client.Read(buf)
		if err != nil {
			return // client disconnected (EOF) or error
		}
		client.Write(buf[:n])
	}
}

func HandleDelimiter(client net.Conn) {
	defer client.Close()
	scanner := bufio.NewScanner(client)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i, b := range data {
			if b == '|' {
				return i + 1, data[:i], nil
			}
		}

		if atEOF {
			if len(data) > 0 {
				fmt.Println("incomplete data in buffer:", string(data))
				return 0, nil, fmt.Errorf("incomplete message: missing |")
			}
			return 0, nil, nil
		}
		return 0, nil, nil
	})

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("connection closed with error:", err)
	}
}

func LengthPayload(client net.Conn) {
	defer client.Close()
	dbPath := getDatabasePath()

	for {
		buf := make([]byte, 1)
		_, err := io.ReadFull(client, buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected cleanly")
			} else {
				fmt.Println("client disconnected:", err)
			}
			return
		}
		fmt.Println("msg type", buf[0])

		buf = make([]byte, 4)
		_, err = io.ReadFull(client, buf)
		if err != nil {
			fmt.Println("payload length reading wasn't successful---------\n", err)
			return
		}

		payloadLength := binary.BigEndian.Uint32(buf) // length in binary int
		fmt.Println("payload length:", payloadLength)

		buf = make([]byte, payloadLength)
		_, err = io.ReadFull(client, buf)
		if err != nil {
			fmt.Println("invalid query ---------\n", err)
			return
		}
		fmt.Println(string(buf))

		database, err := mydb.Open(dbPath)
		if err != nil {
			fmt.Println("database error:", err)
			continue
		}

		rows, err := mydb.Query(database, string(buf))
		if err != nil {
			fmt.Println("query error:", err)
			database.Close()
			continue
		}

		if rows == nil {
			fmt.Println("rows is nil")
			database.Close()
			continue
		}

		// get column names first
		cols, _ := rows.Columns()
		fmt.Println("Columns:", cols)

		// make a slice of any, one per column
		dest := make([]any, len(cols))

		// make a slice of pointers into dest
		ptrs := make([]any, len(cols))
		for i := range dest {
			ptrs[i] = &dest[i]
		}

		// scan into the pointers and print results
		rowCount := 0
		for rows.Next() {
			fmt.Println("----------------------------------------------")
			rows.Scan(ptrs...)

			// Print column names and values
			for i, col := range cols {
				fmt.Printf("%s: %v | ", col, dest[i])
			}
			fmt.Println()
			rowCount++
		}

		if rowCount == 0 {
			fmt.Println("No rows returned from query")
		} else {
			fmt.Printf("✓ Query returned %d row(s)\n", rowCount)
		}

		rows.Close()
		database.Close()
	}
}
