# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o chatApp ./

RUN chmod +x /app/chatApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/chatApp /app

CMD [ "/app/chatApp" ]