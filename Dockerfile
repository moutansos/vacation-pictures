FROM golang:1.23.1-alpine3.20 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main .

FROM alpine:3.14.2

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/vacations.json .
COPY --from=builder /app/pages ./pages
COPY --from=builder /app/static ./static

CMD ["./main"]
