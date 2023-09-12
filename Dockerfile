FROM golang:1.21
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=8088
WORKDIR /golangIM

COPY . .

WORKDIR /golangIM
RUN go build -o main .
EXPOSE 8088

RUN chmod +x main
ENTRYPOINT ["./main"]
