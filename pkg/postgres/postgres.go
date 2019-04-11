package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/kylegrantlucas/platform-exercise/models"
	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	Connection *sql.DB
}

var DB *DatabaseConnection

func CreateDatabase(host, port, user, password, dbName string) (*DatabaseConnection, error) {
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

	return &DatabaseConnection{Connection: db}, nil
}

func (d *DatabaseConnection) CreateUser(email, name, password string) models.User {
	return models.User{}
}

func (d *DatabaseConnection) UpdateUserByUUID(email, name, password string) models.User {
	return models.User{}
}

func (d *DatabaseConnection) GetUserByEmail(email string) models.User {
	return models.User{}
}

func (d *DatabaseConnection) SoftDeleteUserByUUID(email string) models.User {
	return models.User{}
}

func (d *DatabaseConnection) CreateSession(userUUID string, expiresAt time.Time) models.Session {
	return models.Session{}
}

func (d *DatabaseConnection) GetSessionByUserUUID(userUUID string) models.Session {
	return models.Session{}
}

var queries = map[string]string{
	"create_user":              "",
	"create_session":           "",
	"get_user_by_email":        "",
	"update_user_by_uuid":      "",
	"soft_delete_user_by_uuid": "",
	"get_session_by_uuid":      "",
}
