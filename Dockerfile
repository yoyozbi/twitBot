FROM golang:alpine3.15 as builder

ARG opts

COPY . /app
WORKDIR /app

RUN env ${opts} go build -o /app/main src/main.go

FROM golang:alpine3.15 as runner

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]