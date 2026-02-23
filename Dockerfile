# generate binary file for deploy container
FROM golang:1.25.5-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# deploy container
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# auto refresh env for local dev env
FROM golang:1.25.5 as dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD ["air"]
