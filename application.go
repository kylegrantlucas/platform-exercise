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
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func attachHandlers(router *mux.Router) {
	// User Handlers
	router.HandleFunc("/users", user.Create).Methods("POST")
	router.HandleFunc("/users", user.Delete).Methods("DELETE")
	router.HandleFunc("/users", user.Update).Methods("PUT")

	// Session Handlers
	router.HandleFunc("/sessions", session.Create).Methods("POST")
	router.HandleFunc("/sessions", session.Delete).Methods("DELETE")
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
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), n)
	if err != nil {
		log.Fatal(err)
	}
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
