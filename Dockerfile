# build image
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/air-verse/air@latest

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .
RUN /go/bin/sqlc generate

# runtime image
FROM golang:1.25.3-alpine
WORKDIR /app

COPY --from=builder /go/bin/air /go/bin/air
COPY --from=builder /app /app

EXPOSE 3000
CMD ["air", "-c", ".air.toml"]
