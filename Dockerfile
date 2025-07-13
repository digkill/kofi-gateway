FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN goose -dir ./migrations mysql "$MYSQL_DSN" up || true

RUN go build -o gateway ./cmd/main.go

CMD [ "./gateway" ]
