# SQLite Wire Protocol

A custom TCP wire protocol for SQLite, written in Go. Lets clients connect over TCP and send SQL queries to a remote SQLite database — similar to how you'd connect to Postgres, but backed by SQLite.

## How it works

The protocol uses a simple binary framing format:

```
[1 byte: message type][4 bytes: payload length (big-endian uint32)][N bytes: payload]
```

### Message Types

| Byte   | Type    | Description              |
|--------|---------|--------------------------|
| `0x01` | Auth    | Authentication handshake |
| `0x02` | Query   | SQL query string         |

### Response Format

```
[0x01] header packet  — column names
[0x02] row packets    — one per row
[0x03] done signal    — end of result set
```

## Project Structure

```
tcp/
├── server/       — standalone server binary
│   ├── main.go
│   ├── handler/  — protocol parsing, response encoding
│   └── db/       — SQLite integration
└── driver/       — client driver (importable package)
    └── driver.go
```

## Starting the Server

```bash
cd server
go run . <port>
```

Example:
```bash
go run . 8000
```

Server listens on `:<port>` and accepts TCP connections.

## Using the Driver

```go
import "github.com/pandaAritra/tcp/driver"

client := driver.NewQrer()
err := client.Connect("localhost:8000")
if err != nil {
    log.Fatal(err)
}

client.QueryString("SELECT * FROM users")
```

## Testing with netcat

```bash
nc <server-ip> 8000
```

Note: binary protocol messages can't be sent manually via netcat — use the driver package or write a Go client.

## Status

- [x] TCP server with concurrent client handling
- [x] Binary length-prefixed protocol
- [x] Client driver package
- [ ] SQLite query execution
- [ ] Auth handshake
- [ ] Response encoding
