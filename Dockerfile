FROM golang:1.27.3-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates make && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY scripts/install-margo.sh scripts/install-margo.sh
RUN sh scripts/install-margo.sh
COPY . .
RUN make check

FROM nginx:alpine
COPY --from=builder /src/public /usr/share/nginx/html
EXPOSE 80
