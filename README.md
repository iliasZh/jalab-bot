# How to run

`cd` to the repository's root.

1. Run `docker-compose up -d`.
2. Create the schema in the launched PostgreSQL. If you have `psql`, run `psql -U admin -d jalab_bot -a -f ./jalab_bot-schema.sql`, then type `admin` for password.
3. Run `go run cmd/main.go -bot_token <your-bot-token>`.
