package config

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func ConnectToDatabase(databaseDNS string) (*sql.DB, error) {
	return sql.Open("sqlite", databaseDNS)
}
