[![Package build](https://github.com/NDOY3M4N/api-calculator/actions/workflows/deploy.yml/badge.svg?event=release)](https://github.com/NDOY3M4N/api-calculator/actions/workflows/deploy.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/NDOY3M4N/api-calculator)
![Golang Version](https://img.shields.io/badge/Go-1.23-blue?logo=go)

# Backend API - Calculator

The goal of this project is to create an http+json API for a calculator service.

> [!WARNING]
> This project is just for learning purposes and is not intended to be used in production.

## Endpoints

- `/api/v1/login` - Log the user
- `/api/v1/sum` - Sum two numbers
- `/api/v1/add` - Add two numbers
- `/api/v1/substract` - Substract two numbers
- `/api/v1/multiply` - Multiply two numbers
- `/api/v1/divide` - Divide two numbers

## Usage

To use the Docker package, pull the image from the GitHub registry and run it:

```sh
docker pull ghcr.io/NDOY3M4N/api-calculator:latest
docker compose up
```

> [!TIP]
> Make sure to use the `./docker-compose.yml` file in the root of the repository when running the container.

Before you can perform any operations you'll first need to login

```bash
curl -X POST http://localhost:3000/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"pseudo":"b4tm4n"}'
```

Now you can perform any operation using the token that you received.

```bash
curl -X POST http://localhost:3000/api/v1/add \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_SECRET_TOKEN' \
  -d '{"number1":2, "number2": 2}'
```

## Overview

With this API, you can:

- Perform arithmetic operations (addition, subtraction, multiplication, division) on two numbers
- Authenticate with a token to access the API
- Store calculation history in a database
- Handle floating-point numbers
- Benefit from rate limiting to prevent API abuse
- Receive a unique request ID for each request

## Additional Tasks

- [ ] Create an associated http client that can work with the calculator API.
- [ ] Create a frontend that makes use of your API.
- [x] Embed or create the sqlite db in a temp file
- [x] Add in token authentication to prevent anyone unauthorized from using the API
- [x] Add in a database to keep track of all of the calculations that have taken place
- [x] Add in support for floating point numbers as well.
- [x] Add in rate limiter to prevent misuse of the API
- [x] Add in a middleware that adds a request ID to the http.Request object.
