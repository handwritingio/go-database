FROM golang:1.5

RUN go get -v \
  github.com/golang/lint/golint \
  golang.org/x/tools/cmd/godoc \
  golang.org/x/tools/cmd/vet \
  github.com/handwritingio/waiter

ENV SRC_PATH /go/src/github.com/handwritingio/go-database
RUN mkdir -p $SRC_PATH
WORKDIR $SRC_PATH

ADD . $SRC_PATH

# for godoc
EXPOSE 9000

RUN go get -v -d ./...
RUN go install -v ./...
