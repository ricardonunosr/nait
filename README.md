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
docker run -d -p 8080:8080 --rm --name nait_app nait
```