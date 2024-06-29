FROM golang:1.22-alpine

RUN apk update && apk add --no-cache git

WORKDIR /opt/go-app

COPY . .

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]