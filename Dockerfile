FROM golang:buster as builder

COPY . /app
WORKDIR /app

RUN go build -o /app/main src/main.go

FROM golang:buster as runner

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]