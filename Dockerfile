FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN go get github.com/jackc/pgconn@v1.13.0
RUN go get github.com/jackc/pgconn@v1.13.0
RUN go get github.com/hashicorp/serf/serf@v0.9.8
RUN go get go.opentelemetry.io/otel/exporters/metric/prometheus@v0.15.0

WORKDIR /app/bin

RUN CGO_ENABLED=0 GOOS=linux go build -o kvs ../cmd
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
  wget -qO/go/bin/grpc_health_probe \
  https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /go/bin/grpc_health_probe


# Run container
# FROM scratch
# 
# COPY --from=builder /app/bin/kvs .
# 
# EXPOSE 8080

FROM alpine:latest

COPY --from=builder /app/bin/kvs .
COPY --from=builder /go/bin/grpc_health_probe /bin/grpc_health_probe

# for http
# EXPOSE 8080

# CMD ["/kvs", "-m", "true", "-d", "true"]
# CMD ["/kvs", "-t", "true", "-m", "true", "-d", "true"]
# CMD ["/kvs", "-t", "true", "-m", "true"]
# CMD ["/kvs", "-t", "true"]
ENTRYPOINT ["/kvs"]



