# Fender Digital Platform Engineering Challenge

## Description

Design and implement a RESTful web service to facilitate a user authentication system. The authentication mechanism should be *token based*. Requests and responses should be in **JSON**.

## Usage

### Migrate

(These steps will assume there is a local instance of Postgres running, a guide for Postgres on MacOS can be found [here](https://www.codementor.io/engineerapart/getting-started-with-postgresql-on-mac-osx-are8jcopb#3-configuring-postgres))

This application uses [golang-migrate](https://github.com/golang-migrate/migrate) to handle running migrations, you can find install instructions [here](https://github.com/golang-migrate/migrate/tree/master/cli)

`psql -c 'create database platform_exercise'`

`migrate -source file://db/migrations -database 'postgres://localhost:5432/platform_exercise?sslmode=disable' up`

### Run

`env PG_HOST=localhost PG_DB_NAME=platform_exercise PG_PORT=5432 JWT_KEY=fenderdigital go run application.go`

### Example Queries

### Test

`go test ./...`

## Thoughts

### Security

#### Passwords

In order to ensure security of passwords at rest we use the built in golang bcrypt to hash+salt our passwords. The advantage is that bcrypt provides us with secure salts (generated from crypto/rand) and helper methods for checking an submitted password against the stored hash+salted password.

#### Email Validation

When a user submits an email for registration, we do a quick SMTP check against the address to ensure it is a real valid address, this saves us from having to rely on a regular expression for email validation (https://davidcel.is/posts/stop-validating-email-addresses-with-regex/).

#### JWT

To handle tokening we utilize [JWT](). JWT has the great benefit of encoding token expiry and user information entirely within the token itself, removing the need to store and manage the token directly to track it and allowing a client to call the service without any extra metadata (such as a user UUID). Instead to handle expiration we use user "session" that are checked an authenticated with the token, if there was a need to force a user to get a new token, one would simply have to soft delete the session record. Currently these tokens are encoded with a `JWT_KEY` set as an environment variable, but if we later wanted to add on extra security it supports encoding with RSA keypairs, which could then be stored in a secure credential management format (ex: Vault).

#### Future Enhancements

* Roles
* User Reactivation
* Email Verification
  While it's awesome that emails are SMTP checked, in and ideal world we would also send of an email to the user to ensure that they own the inbox.
* API Key Authentication

## Tools Used

### Language

### Database

### Libraries

* [golang-migrate](https://github.com/golang-migrate/migrate) - A tool for running raw sql migrations.
* negroni - A middleware for HTTP requests, allows us to log each request handler easily and recover from panics without crashing
* logrus - A better logger for go
* pq - A pure-go postgres driver, used as the backing driver for sql.DB
* checkmail - A library for SMTP validating email existence
* jwt - A library providing a full JWT implementation with fun things like parsing to headers

## Requirements

**Models**

The **User** model should have the following properties (at minimum):

1. name
2. email
3. password

You should determine what, *if any*, additional models you will need.

**Endpoints**

All of these endpoints should be written from a user's perspective.

1. **User** Registration
2. Login (*token based*) - should return a token, given *valid* credentials
3. Logout - logs a user out
4. Update a **User**'s Information
5. Delete a **User**

**README**

Please include:
- a readme file that explains your thinking
- how to setup and run the project
- if you chose to use a database, include instructions on how to set that up
- if you have tests, include instructions on how to run them
- a description of what enhancements you might make if you had more time.

**Additional Info**

- We expect this project to take a few hours to complete
- You can use Rails/Sinatra, Python, Go, node.js or shiny-new-framework X, as long as you tell us why you chose it and how it was a good fit for the challenge.
- Feel free to use whichever database you'd like; we suggest Postgres.
- Bonus points for security, specs, etc.
- Do as little or as much as you like.

Please fork this repo and commit your code into that fork.  Show your work and process through those commits.
