FROM golang:1.22

RUN mkdir -p "/app/bin"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod vendor
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./cmd/main_entry_point

EXPOSE 8080

CMD ["/app/bin/service"]