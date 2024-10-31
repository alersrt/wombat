package dao

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
)

type PostgreSQLManager struct {
}

func (receiver *PostgreSQLManager) test() {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)

	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name())
	}
}
