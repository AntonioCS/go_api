FROM golang:1.18.1-alpine

WORKDIR /app

COPY ../../project/ .

RUN go mod download

RUN go build -o bin/app cmd/app/main.go

EXPOSE 3000

CMD ["./bin/app"]
