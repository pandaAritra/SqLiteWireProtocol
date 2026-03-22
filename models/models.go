package models

// Request
type Request struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args"`
}

// Response
type Response struct {
	Columns []string `json:"columns,omitempty"`
	Rows    [][]any  `json:"rows,omitempty"`
	Error   string   `json:"error,omitempty"`
}
