
ifeq ($(MAKECMDGOALS),run)
include .env/.run.env
export
endif
run:
	go run cmd/main.go

ifeq ($(MAKECMDGOALS),run-arkady)
include .env/.run-arkady.env
export
endif
run-arkady:
	go run cmd/main.go