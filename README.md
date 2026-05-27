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
| `0x02` | Query   | SQL query string         |

### Response Format (WIP)

Currently returns results directly to stdout. Planned response format:

```
[0x01] header packet  — column names
[0x02] row packets    — one per row
[0x03] done signal    — end of result set
```

## Project Structure

```
SqLiteWireProtocol/
├── cmd/
│   └── main.go              — Server entry point, accepts CLI port arg
├── Handlers/
│   └── ClientHandler.go     — Protocol parsing, query execution, result handling
├── db/
│   ├── db.go                — SQLite wrapper (mydb package)
│   └── test.db              — SQLite database file
├── Utils/
│   ├── bindPort.go          — TCP listener setup
│   ├── getport.go           — CLI port parsing
│   └── makepackage.go       — Package utilities
├── models/
├── go.mod
└── README.md
```

## Starting the Server

```bash
cd SqLiteWireProtocol
go run ./cmd <port>
```

Example:

```bash
go run ./cmd 8000
```

Server listens on `[::]:<port>` (IPv6 + IPv4) and accepts concurrent TCP connections.

## Protocol Usage

Send binary-framed SQL queries:

```
[0x02] [0x00 0x00 0x00 0x29] [SELECT * FROM users where name = "Aritra"]
 ^type          ^length (41 bytes)         ^payload (SQL query)
```

Results are printed to server stdout. Client receives acknowledgment when complete.

## Testing

### Using netcat (informational only)

```bash
nc localhost 8000
```

Note: Binary protocol messages need proper encoding — netcat won't work for actual queries.

### Using Go client

```go
import "net"

conn, _ := net.Dial("tcp", "localhost:8000")
defer conn.Close()

// Send message type (0x02 for query)
conn.Write([]byte{0x02})

// Send payload length as big-endian uint32
payloadLength := uint32(len(query))
binary.Write(conn, binary.BigEndian, payloadLength)

// Send query
conn.Write([]byte(query))
```

## Setup & Requirements

### Database File

1. **Create a fresh database:**
   ```bash
   python3 create_db.py
   ```
   This generates `db/test.db` with sample `users` table.

2. **Update absolute path in code:**
   In `Handlers/ClientHandler.go`, update the database path:
   ```go
   database, err := mydb.Open("/absolute/path/to/db/test.db")
   ```
   (Relative paths are fragile depending on where you run the server from.)

### Dependencies

```bash
go get github.com/mattn/go-sqlite3  # If using CGo SQLite driver
```

Or if using `modernc.org/sqlite` (pure Go):

```bash
go get modernc.org/sqlite
```

## Known Issues & Lessons

- **Relative paths break easily** — Always use absolute paths for database files
- **Database must be opened once** — Currently opens on every query (inefficient). Should open once in `main()` and pass to handlers
- **Response encoding incomplete** — Currently prints results to stdout, not binary-encoded responses
- **No error responses** — Errors logged to stdout only, client receives no feedback

## Status

- [x] TCP server with concurrent client handling
- [x] Binary length-prefixed protocol parsing
- [x] SQLite query execution
- [x] Result scanning and display
- [ ] Binary response encoding (0x01/0x02/0x03)
- [ ] Auth handshake
- [ ] Client driver package
- [ ] Connection pooling
- [ ] Prepared statements

## Next Steps

1. **Refactor database management** — Open once at startup, pass to handlers
2. **Implement response encoding** — Send results back to client, not stdout
3. **Add error handling** — Send error messages to client
4. **Build client driver** — Go package to simplify client code
5. **Add prepared statements** — For better performance and security

## Running Example

```bash
# Terminal 1: Start server
go run ./cmd 8000
# Output: listening on [::]:8000

# Terminal 2: Send query (using Go test client or similar)
# Send: [0x02][length][SELECT * FROM users]

# Terminal 1 shows:
# client is 127.0.0.1:12345
# msg type 2
# payload length: 41
# SELECT * FROM users where name = "Aritra"
# Columns: [id name email age]
# id: 1 | name: Aritra | email: aritra@example.com | age: 24 |
# ✓ Query returned 1 row(s)
```