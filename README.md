# bots-builder

To run bot:

1. Copy ./env/template.env to ./env/.run.env
2. Fill vars values
3. Run `make run`

# Bot ENV variables

- TOKEN - Telegram Bot API token. Required.
- QNR_ID - questionnaire ID. Required.
- DATABASE_HOST - Database host. Required in order to run different servers for different bots.
- DATABASE_PORT - Database port. Sane default in go.Dockerfile.
- DATABASE_NAME - Database name. Sane default in go.Dockerfile.
- DATABASE_USER - Database username. Sane default in go.Dockerfile.
- DATABASE_PASS - Username passwork. Sane default in go.Dockerfile.
- DATABASE_CONNECT_TIMEOUT - Database connect timeout. Sane default in go.Dockerfile.

# Build and Run

Each questionnaire bot declared in an individual Makefile

```
make -f epilepsy.Makefile postgres-start
make -f epilepsy.Makefile docker-build
make -f epilepsy.Makefile bot-start
```

## docker-network-create

Creates docker network

## postgres-start

Starts PostgreSQL with schema and master.

## run-psql

Runs psql client connected to PostgreSQL db.

## docker-build

Builds docker image with the Bot Go application.

## bot-start

Runs the Bot with corresponding environment variables, see --env-file cli argument.

## bot-stop

Runs the Bot.

# Maintenance

## Dump DB and store in the repo

```
pg_dump -U postgres --dbname="bots_builder" --file="./bots_builder-epilepsy-dump.sql"
docker cp pg_bots_builder:/bots_builder-epilepsy-dump.sql ./bots_builder-epilepsy-dump.sql
mv ./bots_builder-epilepsy-dump.sql scripts/migrations/data/epilepsy.sql
```
