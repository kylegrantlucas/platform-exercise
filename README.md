# Fender Digital Platform Engineering Challenge

## Description

Design and implement a RESTful web service to facilitate a user authentication system. The authentication mechanism should be *token based*. Requests and responses should be in **JSON**.

## Usage

### Migrate

(These steps will assume there is a local instance of Postgres running, a guide for Postgres on MacOS can be found [here](https://www.codementor.io/engineerapart/getting-started-with-postgresql-on-mac-osx-are8jcopb#3-configuring-postgres))

This application uses [golang-migrate](https://github.com/golang-migrate/migrate) to handle running migrations, you can find install instructions [here](https://github.com/golang-migrate/migrate/tree/master/cli).

```bash
$ psql -c 'create database platform_exercise'
```

```bash
$ migrate -source file://db/migrations -database 'postgres://localhost:5432/platform_exercise?sslmode=disable' up
```

### Run

```bash
$ env PG_HOST=localhost PG_DB_NAME=platform_exercise PG_PORT=5432 JWT_KEY=fenderdigital go run application.go
```

### Example Queries

### Test

```bash
$ go test ./...
```

## Thoughts

### UUID

Both sessions and user accounts are associated with a UUID. This allows us to ensure accuracy when picking records but also allows us a certain level of obfuscation. With UUIDs the token never has to embed user PII in the token directly, and it makes user ID guessing attacks near impossible.

### Passwords

In order to ensure security of passwords at rest we use the built in golang `bcrypt` package to hash+salt our passwords. The advantage is that bcrypt provides us with secure salts (generated from `crypto/rand`) and helper methods for checking an submitted password against the stored hash+salted password.

On top of this we require the password to not be present in the HaveIBeenPwned database, making it much hard to crack the encryption by matching external dictionaries, and eliminating weak passwords.

### Email Validation

When a user submits an email for registration we do a quick regular expression validation of their email address.

### JWT

To handle tokening we utilize [JWT](). JWT has the great benefit of encoding token expiry and user information entirely within the token itself, removing the need to store and manage the token directly to track it and allowing a client to call the service without any extra metadata (such as a user UUID). Instead to handle expiration we use user "session" that are checked an authenticated with the token, if there was a need to force a user to get a new token, one would simply have to soft delete the session record. Currently these tokens are encoded with `HMAC512` and a `JWT_KEY` set as an environment variable, but if we later wanted to add on extra security it supports encoding with RSA keypairs, which could then be stored in a secure credential management format (ex: Vault).

### Protected Endpoint

All actions other than Create User and Create Session are JWT Token protected.

#### Future Enhancements

* Roles

  Currently we only have one role for all users, we could enable permissions on different levels by implementing roles, and JWT tokens will allow us to communicate roles through the token. These can then be stored against the Users table.
* User Reactivation

  Currently once a user account is deleted, that's it. That email can no longer be used. It would be fairly trivial to allow users to reactivate their account at a later date.
* Email Verification

  While it's awesome that emails are validated, in and ideal world we would also send of an email to the user to ensure that they own the inbox.
* API Key Authentication

  Currently registration is open to anyone who would like to POST at it. You could limit this by implmenting an API token system, where users of the system have to register before they can make calls to the API.
* Rate Limits

  There are no rate limits on the API right now, and in order to perform a 401 with a outdated token we need to do at least 1 database call - in theory this could be abused. Some safe rate limits for reasonable usage would revent this attack vector.
* RSA JWT Encryption

  We could generate an RSA public/private keypair and load the pair into the application via a secret management service to make the tokens more secure.

## Tools Used

### Language

I chose to write this service in Golang. It's personally my most used language, and I gets great performance out-of-the-box with minimal fussing around with external packages and frameworks.

### Database

This application utilizes Postgres. This comes to great benefit as we can use built in constraints to ensure record completeness, and built in extensions for generating compliant UUIDs.

### Libraries

* [golang-migrate](https://github.com/golang-migrate/migrate) - A tool for running raw sql migrations.
* [mux](https://github.com/gorilla/mux) - A router with extensions for Go
* [negroni](https://github.com/urfave/negroni) - A middleware for HTTP requests, allows us to log each request handler easily and recover from panics without crashing
* [logrus](https://github.com/sirupsen/logrus) - A better logger for go
* [pq](https://github.com/lib/pq) - A pure-go postgres driver, used as the backing driver for sql.DB
* [jwt](https://github.com/pascaldekloe/jwt/) - A library providing a full JWT implementation with fun things like parsing to headers
* [pwned-passwords](https://github.com/mattevans/pwned-passwords) - A go library for checking against the HaveIBeenPwned database

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
