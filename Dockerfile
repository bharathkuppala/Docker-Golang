FROM golang:latest

WORKDIR /go/src/coding-prep/docker-go-server

COPY ./ ./

RUN go get -d  -v

RUN go build -o main .

CMD [ "./main" ]