# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o accountApp ./internal/api/cmd

RUN chmod +x /app/accountApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/accountApp /app

CMD [ "/app/accountApp" ]