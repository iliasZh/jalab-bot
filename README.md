# How to run

`cd` to the repository's root.

1. Run `docker-compose up -d`
2. Create the schema in your PostgreSQL using the `...-schema.sql` file in the repository's root
3. Run `go run cmd/main.go -bot_token <your-bot-token>`
