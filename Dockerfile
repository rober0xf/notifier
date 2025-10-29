FROM golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
