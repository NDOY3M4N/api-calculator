default: run-with-docs

run:
	@go run *.go

# Assuming you have `swag` installed on your machine
docs:
	@swag init --v3.1

run-with-docs:
	@make docs
	@make run

# Assuming you have `hyperfine` and `http` installed on your machine
bench-rate:
	@hyperfine --runs 5 "http :3000/api/v1/add number1:=2 number2:=2 --ignore-stdin" --show-output

.PHONY: docs test build 
