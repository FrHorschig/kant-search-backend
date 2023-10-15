package util

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func InitDbConnection() *sql.DB {
	connStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
