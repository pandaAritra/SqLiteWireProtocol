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
│   ├── main.go              — Server entry point, accepts CLI port arg
│   └── client/
│       └── main.go          — Client CLI entry point (REPL & single query)
├── Handlers/
│   └── ClientHandler.go     — Protocol parsing, query execution, result handling
├── db/
│   ├── db.go                — SQLite wrapper (mydb package)
│   └── test.db              — SQLite database file
├── Utils/
│   ├── bindPort.go          — TCP listener setup
│   ├── getport.go           — CLI port parsing
│   └── makepacket.go        — Package utilities (packet encoding)
├── models/
│   └── models.go            — JSON-structured request/response models
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

## Running the Client

We have built a dedicated Go client that supports both an **interactive REPL shell** and a **single-query execution** mode:

### 1. Interactive REPL Mode
To open an interactive session and run multiple queries:
```bash
go run ./cmd/client/main.go <port>
```
Example:
```bash
go run ./cmd/client/main.go 8000
```
This starts an interactive shell where you can type queries directly:
```
sqlite-wire> SELECT * FROM users;
```

### 2. Single-Query Mode
To run a query directly and exit:
```bash
go run ./cmd/client/main.go <port> "<sql_query>"
```
Example:
```bash
go run ./cmd/client/main.go 8000 "SELECT * FROM users;"
```

## Setup & Requirements

### Database File

The server dynamically resolves the database location in the following order:
1. Environment variable `SQLITE_DB_PATH`
2. Local workspace `./db/test.db`
3. Absolute path `/home/panda/Projects/SqLiteWireProtocol/db/test.db`
4. Absolute path `/home/panda/Documents/code/test/SqLiteWireProtocol/db/test.db`

This ensures that the database is detected and works automatically without needing manual code changes.

### Dependencies

```bash
go get modernc.org/sqlite
```

## Status

- [x] TCP server with concurrent client handling
- [x] Binary length-prefixed protocol parsing
- [x] SQLite query execution
- [x] Result scanning and display
- [x] Interactive REPL / Single-Query Client CLI
- [x] Dynamic database path resolution
- [ ] Binary response encoding (0x01/0x02/0x03)
- [ ] Auth handshake
- [ ] Connection pooling
- [ ] Prepared statements