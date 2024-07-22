FROM golang:1.22 AS builder

RUN mkdir -p "/app/bin"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./cmd/main_entry_point

FROM golang:1.22

WORKDIR /app
COPY --from=builder /app/bin/service /app/service
COPY .env.local .env

EXPOSE 8080

CMD ["/app/service"]