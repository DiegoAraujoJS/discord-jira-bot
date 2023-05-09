FROM golang:latest
ADD src /src
WORKDIR /src
RUN go build
CMD ["./go-bot"]
