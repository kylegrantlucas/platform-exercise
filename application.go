package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kylegrantlucas/platform-exercise/handlers/session"
	"github.com/kylegrantlucas/platform-exercise/handlers/user"
	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
	"github.com/pascaldekloe/jwt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func attachHandlers(router *mux.Router) {
	keys := &jwt.KeyRegister{Secrets: [][]byte{[]byte(os.Getenv("JWT_KEY"))}}
	headers := map[string]string{
		"sub": "X-Verified-User-Uuid",
		"sid": "X-Verified-Session-Uuid",
	}

	router.Use(jsonMiddleware)

	// User Handlers
	router.HandleFunc("/users", user.Create).Methods("POST")
	router.Handle("/users", &jwt.Handler{Target: http.HandlerFunc(user.Delete), HeaderBinding: headers, Keys: keys}).Methods("DELETE")
	router.Handle("/users", &jwt.Handler{Target: http.HandlerFunc(user.Update), HeaderBinding: headers, Keys: keys}).Methods("PUT")

	// Session Handlers
	router.HandleFunc("/sessions", session.Create).Methods("POST")
	router.Handle("/sessions", &jwt.Handler{Target: http.HandlerFunc(session.Delete), HeaderBinding: headers, Keys: keys}).Methods("DELETE")
}

func main() {
	// Setup Postgres connection early, so we can fail fast if it doesn't work
	var err error
	postgres.DB, err = postgres.CreateDatabase(os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	if err != nil {
		log.Fatalf("couldn't connect to postgres: %v", err)
	}

	// Setup our mux router, handlers, negroni middleware and logger
	router, n, recovery := mux.NewRouter().StrictSlash(true), negroni.New(), negroni.NewRecovery()
	setupLogger(recovery)
	attachHandlers(router)
	n.Use(recovery)
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	// Setup and startup our HTTP server
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Printf("now serving traffic on port :%v", port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), n)
	if err != nil {
		log.Fatal(err)
	}
}

// Ensures all requests are Content-Type application/json
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func setupLogger(recovery *negroni.Recovery) {
	// Create the logrus logger and set a level
	logger := logrus.New()
	w := logger.WriterLevel(logrus.ErrorLevel)

	// Setups up pretty line logging with stack traces on anything other than a 'production' environment
	if os.Getenv("GO_ENV") != "production" {
		logger.Level = logrus.InfoLevel
		logger.Formatter = &logrus.TextFormatter{ForceColors: true}
		recovery.PrintStack = true
		recovery.Logger = log.New(w, "", 0)
	} else {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	// Output to logrus, add line numbers so we can find our logging statements easier
	log.SetOutput(logger.Writer())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
