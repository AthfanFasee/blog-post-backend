FROM postgres:12-alpine
COPY install-extension.sql /docker-entrypoint-initdb.d
ENV POSTGRES_USER=root \
    POSTGRES_PASSWORD=secret \
    POSTGRES_DB=blogpost
