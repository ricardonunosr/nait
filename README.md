# Nait

## Tech Stack
- [htmx](https://htmx.org)
- [TailwindCSS](https://tailwindcss.com/)
- [Django Templates]()
- [Stripe]()
- [Go]()
- [Fiber]()
- [migrate](https://github.com/golang-migrate/migrate)

## Local Development

```bash
# This will run watch if files are changed and build and restart the server
# You need air for local development:
#      - go install github.com/cosmtrek/air@latest
#      - install tailwindcss CLI
make run
```

# Dockerfile
```bash
docker build . -t nait
docker run -p 80:80 --env-file .env --rm --name nait_app nait
```

# `.env` example

```.env
CLUB_NAME=LimeSet
PORT=80

# Supabase
SUPABASE_URL=
SUPABASE_KEY=

# Database
NAIT_DB_DSN=sqlite3://database.db

# Stripe
STRIPE_KEY=
```