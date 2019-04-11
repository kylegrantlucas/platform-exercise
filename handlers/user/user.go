package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
)

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

	// Check email format validation
	err = checkmail.ValidateFormat(parsedBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Email is invalid, %v"}`, err)))
		return
	}

	newUser, err := postgres.DB.CreateUser(parsedBody.Email, parsedBody.Name, parsedBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	response, err := json.Marshal(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	user, err := postgres.DB.SoftDeleteUserByUUID(r.Header["X-Verified-User-Uuid"][0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
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

func Update(w http.ResponseWriter, r *http.Request) {
	parsedBody, err := parseUserRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := postgres.DB.UpdateUserByUUID(r.Header["X-Verified-User-Uuid"][0], parsedBody.Email, parsedBody.Name, parsedBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
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
