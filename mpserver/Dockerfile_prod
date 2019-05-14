FROM golang:latest
RUN mkdir /go/src/app
ADD . /go/src/app/
WORKDIR /go/src/app
RUN go get ./...
RUN go build -o main .
EXPOSE 8081
CMD ["/go/src/app/main"]
