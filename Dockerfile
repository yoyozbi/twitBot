FROM golang:alpine3.15 as builder

COPY . /app
WORKDIR /app

RUN env ${ops} go build -o /app/main src/main.go

FROM golang:alpine3.15 as runner

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]