package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/kylegrantlucas/platform-exercise/models"
	"github.com/kylegrantlucas/platform-exercise/pkg/password"
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

func (d *DatabaseConnection) CreateUser(email, name, plaintextPassword string) (models.User, error) {
	user := models.User{}

	// Hash + Salt our password prior to creating the user record
	encryptedPassword, err := password.HashAndSalt(plaintextPassword)
	if err != nil {
		return user, err
	}

	// Insert the record
	rows, err := d.Connection.Query(queries["create_user"], email, name, encryptedPassword, time.Now())
	if err != nil {
		return user, err
	}

	// Scan off the result to return to the client
	for rows.Next() {
		err := rows.Scan(&user.UUID, &user.Email, &user.Name, &user.CreatedAt)
		if err != nil {
			return user, err
		}
	}

	// Check to make sure there were no errors during scan
	err = rows.Err()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *DatabaseConnection) UpdateUserByUUID(uuid, email, name, plaintextPassword string) (models.User, error) {
	user := models.User{}
	var rows *sql.Rows

	queryBody := ""
	args := []interface{}{}
	argCount := 0

	if email != "" {
		argCount++
		queryBody += fmt.Sprintf("email=$%v, ", argCount)
		args = append(args, email)
	}

	if name != "" {
		argCount++
		queryBody += fmt.Sprintf("name=$%v, ", argCount)
		args = append(args, name)
	}

	if plaintextPassword != "" {
		argCount++

		encryptedPassword, err := password.HashAndSalt(plaintextPassword)
		if err != nil {
			return user, err
		}

		queryBody += fmt.Sprintf("password=$%v, ", argCount)
		args = append(args, encryptedPassword)
	}

	argCount++
	queryBody += fmt.Sprintf(" where uuid=$%v", argCount)
	args = append(args, uuid)

	// Update the record
	rows, err := d.Connection.Query(fmt.Sprintf(queries["update_user"], queryBody), args...)
	if err != nil {
		return user, err
	}

	// Scan off the result to return to the client
	for rows.Next() {
		err := rows.Scan(&user.UUID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return user, err
		}
	}

	// Check to make sure there were no errors during scan
	err = rows.Err()
	if err != nil {
		return user, err
	}

	return user, nil
}

// func (d *DatabaseConnection) GetUserByUUID(uuid string) (models.User, error) {
// 	return models.User{}, nil
// }

func (d *DatabaseConnection) SoftDeleteUserByUUID(uuid string) (models.User, error) {
	user := models.User{}

	// Soft delete the record
	rows, err := d.Connection.Query(queries["delete_user"], time.Now(), uuid)
	if err != nil {
		return user, err
	}

	// Scan off the result to return to the client
	for rows.Next() {
		err := rows.Scan(&user.UUID, &user.Email, &user.Name, &user.CreatedAt, &user.DeletedAt, &user.UpdatedAt)
		if err != nil {
			return user, err
		}
	}

	// Check to make sure there were no errors during scan
	err = rows.Err()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *DatabaseConnection) CreateSession(userUUID string, expiresAt time.Time) (models.Session, error) {
	return models.Session{}, nil
}

func (d *DatabaseConnection) GetSessionByUserUUID(userUUID string) (models.Session, error) {
	return models.Session{}, nil
}

var queries = map[string]string{
	"create_user":                      "insert into users (email, name, password, created_at) values ($1, $2, $3, $4) returning uuid, email, name, created_at;",
	"create_session":                   "insert into sessions (uuid, created_at, expires_at) values ($1, $2, $3) returning uuid, created_at, expires_at;",
	"update_user_by_uuid":              "update users set %v returning uuid, email, name, created_at, updated_at;",
	"soft_delete_user_by_uuid":         "update users set deleted_at=$1, updated_at=$1 where uuid=$2 returning uuid, email, name, created_at, updated_at, deleted_at;",
	"soft_delete_session_by_user_uuid": "update sessions set deleted_at=$1 where user_uuid=$2 returning user_uuid, created_at, expires_at, deleted_at;",
	"get_session_by_user_uuid":         "select user_uuid, created_at, expires_at WHERE expires_at > $1 AND deleted_at IS NULL LIMIT 1;",
	// "get_user_by_uuid":        "",
}
