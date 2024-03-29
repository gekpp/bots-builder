
.PHONY: start-postgres
start-postgres:
	docker run --rm -d --name pg_bots_builder \
	-e POSTGRES_PASSWORD=admin \
	-e POSTGRES_DB=bots_builder \
	-e POSTGRES_USER=postgres \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v `pwd`/scripts/migrations:/docker-entrypoint-initdb.d/ \
	-p 5432:5432 \
	postgres:15.1




.PHONY: run
ifeq ($(MAKECMDGOALS),run)
include .env/.run.env
export
endif
run:
	go run cmd/*.go

.PHONY: run-arkady
ifeq ($(MAKECMDGOALS),run-arkady)
include .env/.run-arkady.env
export
endif
run-arkady:
	go run cmd/*.go

.PHONY: run-anton
ifeq ($(MAKECMDGOALS),run-anton)
include .env/.run-anton.env
export
endif
run-anton:
	go run cmd/*.go