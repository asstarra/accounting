package db

import (
	"fmt"
)

// Переменные для подключения к БД.
var DB struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

// Получение строки подключения.
func DataSourseTcp() string {
	return fmt.Sprint(DB.Host, ":", DB.Password, "@tcp/", DB.Database)
}
