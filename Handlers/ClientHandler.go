package handlers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"unicode/utf8"

	utils "github.com/pandaAritra/sqliteWireProtocol/Utils"
	mydb "github.com/pandaAritra/sqliteWireProtocol/db"
	"github.com/pandaAritra/sqliteWireProtocol/models"
)

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
	dbPath := utils.GetDatabasePath()

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
			SendResponse(client, &models.Response{Error: fmt.Sprintf("database error: %v", err)})
			continue
		}

		rows, err := mydb.Query(database, string(buf))
		if err != nil {
			fmt.Println("query error:", err)
			SendResponse(client, &models.Response{Error: fmt.Sprintf("query error: %v", err)})
			database.Close()
			continue
		}

		if rows == nil {
			fmt.Println("rows is nil")
			SendResponse(client, &models.Response{Error: "query returned no result set"})
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
		var allRows [][]any
		for rows.Next() {
			fmt.Println("----------------------------------------------")
			if err := rows.Scan(ptrs...); err != nil {
				fmt.Println("scan error:", err)
				break
			}

			// Print column names and values
			for i, col := range cols {
				fmt.Printf("%s: %v | ", col, dest[i])
			}
			fmt.Println()

			// We need to copy dest slice elements, because scanning overwrites them
			rowVals := make([]any, len(cols))
			for i, val := range dest {
				if b, ok := val.([]byte); ok {
					if utf8.Valid(b) {
						rowVals[i] = string(b)
					} else {
						rowVals[i] = b
					}
				} else {
					rowVals[i] = val
				}
			}
			allRows = append(allRows, rowVals)
			rowCount++
		}

		if rowCount == 0 {
			fmt.Println("No rows returned from query")
		} else {
			fmt.Printf("✓ Query returned %d row(s)\n", rowCount)
		}

		resp := &models.Response{
			Columns: cols,
			Rows:    allRows,
		}
		if err := SendResponse(client, resp); err != nil {
			fmt.Println("failed to send response:", err)
		}

		rows.Close()
		database.Close()
	}
}
