default: run-with-docs

run:
	@go run *.go

docs:
	@swag init --v3.1

run-with-docs:
	@make docs
	@make run

.PHONY: docs test build 
