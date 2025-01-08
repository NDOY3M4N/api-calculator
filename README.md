# Backend API - Calculator

The goal of this project is to create an http+json API for a calculator service.

## Overview

This calculator is stateless (for now), meaning that there is no data stored in a database or in memory.

## Additional Tasks

- [ ] Add in token authentication to prevent anyone unauthorized from using the API
- [ ] Add in a database to keep track of all of the calculations that have taken place
- [ ] Add in support for floating point numbers as well.
- [ ] Create an associated http client that can work with the calculator API.
- [ ] Create a frontend that makes use of your API.
- [x] Add in rate limiter to prevent misuse of the API
- [x] Add in a middleware that adds a request ID to the http.Request object.
