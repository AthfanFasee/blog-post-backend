# Build Stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
ARG VERSION
ENV VERSION=$VERSION
RUN GOOS=linux GOARCH=amd64 go build -ldflags='-s -X main.version=$VERSION' -o main ./cmd/api
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY /migrations ./migration

EXPOSE 3001
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]