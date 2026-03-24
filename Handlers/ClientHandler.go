package handlers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	db "github.com/pandaAritra/sqliteWireProtocol/db"
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

func HandelDelimeter(client net.Conn) {
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
			fmt.Println("payload length reading wasnt successful---------\n", err)
			return
		}

		payloadLength := binary.BigEndian.Uint32(buf) // length in binary int
		fmt.Println("payload length:", payloadLength)

		buf = make([]byte, payloadLength)
		_, err = io.ReadFull(client, buf)
		if err != nil {
			fmt.Println("invalid wuarry ---------\n", err)
			return
		}
		fmt.Println(string(buf))

		database, err := db.Open("../db/test.db")

		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query(database, string(buf))

		// get column names first
		cols, _ := rows.Columns()
		fmt.Println(cols)

		// make a slice of any, one per column
		dest := make([]any, len(cols))

		// make a slice of pointers into dest
		ptrs := make([]any, len(cols))
		for i := range dest {
			ptrs[i] = &dest[i]
		}

		// now scan into the pointers
		for rows.Next() {
			fmt.Println("----------------------------------------------")
			rows.Scan(ptrs...)
			fmt.Printf("%T\n", dest)
		}

	}
}
