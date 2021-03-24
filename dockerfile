FROM golang:latest

# Install golang-migrate/migrate for migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /code
COPY . .
RUN go get -u ./...

LABEL maintainer=9r1nc3w1ll@gmail.com
