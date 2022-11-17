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

WORKDIR /app/bin

RUN CGO_ENABLED=0 GOOS=linux go build -o kvs ../cmd


# Run container
FROM scratch

COPY --from=builder /app/bin/kvs .

EXPOSE 8080

# CMD ["/kvs", "-t", "true", "-m", "true"]
CMD ["/kvs"]



