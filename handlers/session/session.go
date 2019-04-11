package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kylegrantlucas/platform-exercise/pkg/password"
	"github.com/pascaldekloe/jwt"

	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
)

func Create(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedBody := sessionRequest{}
	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := postgres.DB.GetUserByEmail(parsedBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	if !password.ComparePlaintextWithEncypted(parsedBody.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentTime := time.Now()
	expireTime := currentTime.Add(24 * time.Hour)
	session, err := postgres.DB.CreateSession(user.UUID, expireTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	var claims jwt.Claims
	claims.Issuer = "fender"
	claims.Subject = user.UUID
	claims.NotBefore = jwt.NewNumericTime(currentTime)
	claims.Issued = jwt.NewNumericTime(currentTime)
	claims.Expires = jwt.NewNumericTime(expireTime)
	claims.Set = map[string]interface{}{"sid": session.UUID}

	token, err := claims.HMACSign(jwt.HS512, []byte(os.Getenv("JWT_KEY")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%v"}`, string(token))))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	_, err := postgres.DB.SoftDeleteSessionByUUID(r.Header["X-Verified-Session-Uuid"][0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "%v"}`, err)))
		return
	}

	w.WriteHeader(http.StatusOK)
}

type sessionRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
