package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
)

func Create(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedBody := UserRequest{}
	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if parsedBody.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Password must be set"}`))
		return
	}

	err = checkmail.ValidateHost(parsedBody.Email)
	if _, ok := err.(checkmail.SmtpError); ok && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Email is invalid"}`))
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

	_, err = w.Write(response)
	if err != nil {
		log.Print("error writing response")
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func Update(w http.ResponseWriter, r *http.Request) {
}

type UserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}
