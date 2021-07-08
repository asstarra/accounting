package db

// Название полей и таблиц.
var Tab map[string]Table

type Table struct {
	Name    string            `json:"Name"`
	Columns map[string]Column `json:"Columns"`
}

type Column struct {
	Name string `json:"Name"`
}
