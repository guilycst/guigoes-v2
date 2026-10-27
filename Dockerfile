FROM python:3.13-slim AS builder
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates make && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY scripts/install-margo.sh scripts/install-margo.sh
RUN sh scripts/install-margo.sh
COPY . .
RUN make check

FROM nginx:alpine
COPY --from=builder /src/public /usr/share/nginx/html
EXPOSE 80
