FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/molefuckgo/gin-blog
COPY . $GOPATH/src/github.com/molefuckgo/gin-blog
RUN go build .

EXPOSE 8080
ENTRYPOINT ["./gin-blog"]