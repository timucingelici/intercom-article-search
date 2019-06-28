FROM golang:latest

ENV PROJECT_DIR /go/src/github.com/timucingelici/intercom-article-search

RUN apt update

RUN mkdir -p $PROJECT_DIR
COPY . $PROJECT_DIR
WORKDIR $PROJECT_DIR
CMD export GOPATH=/go && go build -o $GOPATH/bin/run-app $PROJECT_DIR/main.go && $GOPATH/bin/run-app
