# Build Stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
ARG VERSION
ARG CURRENT_TIME
ENV VERSION=$VERSION \
CURRENT_TIME=$CURRENT_TIME
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -X main.version=$VERSION -X main.buildTime=$CURRENT_TIME" -o main ./cmd/api
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.16
WORKDIR /app
ENV DB_DSN=postgres://root:secret@pstogres:5432/blogpost?sslmode=disable
RUN echo DB_DSN
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY /migrations ./migration

# Import environment variables from the app.env file
RUN set -o allexport; source ./app.env; set +o allexport

EXPOSE 3001
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/main", "--env", "production", "--cors-trusted-origins", "$CORS_TRUSTED_ORIGINS" ]