.PHONY: docker-network-create
docker-network-create:
	@docker network create bots_builder_epilepsy


.PHONY: postgres-start
postgres-start:
	@docker run -d \
	--name epilepsy_pg \
	--restart=always \
	--network=bots_builder_epilepsy \
	-e POSTGRES_PASSWORD=admin \
	-e POSTGRES_DB=bots_builder \
	-e POSTGRES_USER=postgres \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v `pwd`/scripts/migrations/data/epilepsy.sql:/docker-entrypoint-initdb.d/02-data.sql \
	postgres:15.1

.PHONE: run-psql
run-psql:
	@docker run -it --rm \
	--network bots_builder_epilepsy \
	postgres:15.1 \
	psql -h epilepsy_pg -U postgres -d bots_builder

# --dbname="bots_builder" --file="bots_builder-epilepsy-dump.sql" --if-exists

.PHONY: run-anton
ifeq ($(MAKECMDGOALS),run-anton)
include .env/.run-anton.env
export
endif
run-anton:
	go run cmd/*.go

.PHONY: docker-build
docker-build:
	docker build \
		-f scripts/deploy/epilepsy-docker/go.Dockerfile \
		-t epilepsy_bot \
		.

.PHONY: bot-start
bot-start:
	docker run \
		-d \
		--network=bots_builder_epilepsy \
		--restart=always \
		--env-file=.env/.bot-start.env \
		--name=epilepsy_bot \
		epilepsy_bot:latest 

.PHONY: bot-stop
bot-stop:
	docker stop epilepsy_bot && docker rm epilepsy_bot

