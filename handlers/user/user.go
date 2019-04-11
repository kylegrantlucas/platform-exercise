package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
	hibp "github.com/mattevans/pwned-passwords"
)

// Create is a handler that creates a user with the given parameters
func Create(w http.ResponseWriter, r *http.Request) {
	parsedBody, err := parseUserRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check password
	if parsedBody.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Password must be set"}`))
		return
	}

	// Check the password against HaveIBeenPwned
	client := hibp.NewClient()
	pwned, err := client.Pwned.Compromised(parsedBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	if pwned {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Password is in the HaveIBeenPwned database"}`))
		return
	}

	// Check email format validation
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(e) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Email is invalid, %v"}`, err)))
		return
	}

	// Create the new user
	newUser, err := postgres.DB.CreateUser(parsedBody.Email, parsedBody.Name, parsedBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	// Marshal the user for response
	response, err := json.Marshal(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Delete is a handler that deletes a user with the UUID provided in the JWT token
func Delete(w http.ResponseWriter, r *http.Request) {
	_, err := postgres.DB.GetSessionByUUID(r.Header["X-Verified-Session-Uuid"][0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the user
	user, err := postgres.DB.SoftDeleteUserByUUID(r.Header["X-Verified-User-Uuid"][0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	// Marshal the user for response
	response, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Update is a handler that updates a user with the UUID provided in the JWT token with the given parameters
func Update(w http.ResponseWriter, r *http.Request) {
	parsedBody, err := parseUserRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := postgres.DB.GetSessionByUUID(r.Header["X-Verified-Session-Uuid"][0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	if session.DeletedAt != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if parsedBody.Password != "" {
		// Check the password against HaveIBeenPwned
		client := hibp.NewClient()
		pwned, err := client.Pwned.Compromised(parsedBody.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
			return
		}

		if pwned {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Password is in the HaveIBeenPwned database"}`))
			return
		}
	}

	user, err := postgres.DB.UpdateUserByUUID(r.Header["X-Verified-User-Uuid"][0], parsedBody.Email, parsedBody.Name, parsedBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	if user.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func parseUserRequest(r *http.Request) (userRequest, error) {
	parsedBody := userRequest{}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return parsedBody, err
	}

	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		return parsedBody, err
	}

	return parsedBody, nil
}

type userRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}
