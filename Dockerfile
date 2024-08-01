FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-auth-api

RUN chmod +x go-auth-api

EXPOSE 3000

CMD ["./go-auth-api"]