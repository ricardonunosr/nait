# Nait

## Tech Stack
- [htmx](https://htmx.org)
- [TailwindCSS](https://tailwindcss.com)
- [Django Templates](https://docs.djangoproject.com/en/4.2/topics/templates)
- [Stripe](https://stripe.com/docs/api)
- [Go](https://go.dev)
- [Fiber](https://docs.gofiber.io)
- [migrate](https://github.com/golang-migrate/migrate)
- [SQLite](https://www.sqlite.org/index.html)
- [Supabase](https://supabase.com)

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