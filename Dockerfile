FROM golang:1.4

RUN go get -v \
  github.com/golang/lint/golint \
  golang.org/x/tools/cmd/godoc \
  golang.org/x/tools/cmd/vet

ENV SRC_PATH /go/src/github.com/graciouseloise/go-database
RUN mkdir -p $SRC_PATH
WORKDIR $SRC_PATH

ADD . $SRC_PATH

RUN go get -v -d ./...
RUN go install -v ./...
