# bots-builder

To run bot:

1. Copy ./env/template.env to ./env/.run.env
2. Fill vars values
3. Run `make run`

```
pg_dump -U postgres --dbname="bots_builder" --file="./bots_builder-epilepsy-dump.sql"
docker cp pg_bots_builder:/bots_builder-epilepsy-dump.sql ./bots_builder-epilepsy-dump.sql
mv ./bots_builder-epilepsy-dump.sql scripts/migrations/data/epilepsy.sql
```
