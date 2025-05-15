FROM migrate/migrate:v4.18.1

ENV POSTGRES_HOST="localhost"
ENV POSTGRES_PORT="5432"
ENV POSTGRES_USER="postgres"
ENV POSTGRES_PASSWORD="5432"
ENV POSTGRES_DB="postgres"
ENV SSL_MODE="disable"

WORKDIR /migrations

COPY ./schema/migrations /migrations

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["migrate -path /migrations -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${SSL_MODE}\" up"]
