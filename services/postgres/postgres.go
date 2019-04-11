package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateDatabase(host, port, user, password, dbName string) (*sql.DB, error) {
	sqlURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbName,
	}

	// Create the database connection
	db, err := sql.Open("postgres", sqlURL.String()+"?sslmode=disable")
	if err != nil {
		log.Printf("Error opening db connection: %v", err)
		return nil, err
	}

	// Test ping the database to make sure the connection is good
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging db: %v", err)
		return nil, err
	}

	return db, nil
}
