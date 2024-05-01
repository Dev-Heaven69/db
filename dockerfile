FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod tidy

COPY . .

RUN go build -o main internal/cmd/main.go

EXPOSE 5000

CMD ["./main"]

