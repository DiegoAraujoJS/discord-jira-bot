# ./app/Dockerfile
FROM golang:1.17-alpine as builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o bot

FROM alpine:3.14
COPY --from=builder /src/bot /bot
COPY config.json /
RUN chmod +x /bot
CMD ["/bot"]
