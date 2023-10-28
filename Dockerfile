FROM golang:1.21.3-bullseye as build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

RUN apt-get update \
    && apt-get install --no-install-recommends --no-install-suggests -y \
    make

COPY data/*.go ./data/
COPY db/*.go ./db/
COPY handlers/*.go ./handlers/
COPY static/* ./static/
COPY views/ ./views/
COPY main.go ./
COPY .env ./
COPY Makefile ./
COPY tailwind.config.js ./
COPY database.db.seed ./database.db

# Build
RUN make build

FROM alpine:latest AS build-release-stage

WORKDIR /app/dist
EXPOSE 80
EXPOSE 443

COPY --from=build-stage /app/dist /app/dist

# Run
ENTRYPOINT ["/app/dist/server"]