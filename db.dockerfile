FROM postgres:12-alpine
COPY install-extension.sql /docker-entrypoint-initdb.d