
# Go Boilerplate

Just writing some golang with things that I like to see in a Web API

## Features

- Hexagonal Architecture (kinda overengineering but ok, also, just wrote like this to see how it goes)
- Simple routing with chi
- Centralized encoding and decoding
- Centralized error handling
- Versioned HTTP Handler
- SQL type safety with SQLC
- Migrations with golang migrate
- PASETO tokens instead of JWT
- Access and Refresh Tokens
- Tests that uses Testify Testcontainers instead of mocks
- Testing scripts that uses cURL and jq (f* Postman)

## Required dependencies

- jq
- golang-migrate
- docker
- sqlc

## TODO

- Use testcontainers with test suite
- Associate error logs with request
