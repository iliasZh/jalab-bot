# What it does
Nothing really useful.

The most common command `/todaysjalab` picks a group member at random (for the day) and @mentions them.
We used it just to spam random colleagues with @mentions.

The other common, somewhat useful command, `/yaxshi` is for giving thanks to a colleague.
With its help we had our own social credit system :)

"jalab" and "yaxshi" are Uzbek words.
Don't ask why, it's just Uzbek. All other text for this bot is in Russian.

Other commands are for statistics and customizability.

# How to run

`cd` to the repository's root.

1. Run `docker-compose up -d`.
2. Create the schema in PostgreSQL. The schema file is in repo's root. 
If you have `psql`, run `psql -U admin -d jalab_bot -a -f ./jalab_bot-schema.sql`, then type `admin` for password.
3. Run `go run cmd/main.go -bot_token <your-bot-token>`.
