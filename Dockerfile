FROM ubuntu:14.04
RUN apt-get update
RUN apt-get install -y wget mercurial git gcc
RUN apt-get clean -y
RUN wget -qO- https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz | tar vxzC /usr/local
ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:$PATH
RUN mkdir /usr/local/gopath
ENV GOPATH /usr/local/gopath

RUN go get -v github.com/abdulkadiryaman/hrotti
RUN go get -v github.com/hashicorp/consul
#WORKDIR /usr/local/gopath/src/github.com/goskydome/mqtt-stack
#RUN go get -d -v ./... && go build -v ./... && go install

ENV PATH $GOPATH/bin:$PATH

ENTRYPOINT go get -v github.com/goskydome/mqtt-stack && cd /usr/local/gopath/src/github.com/goskydome/mqtt-stack && go get -d -v ./... && go build -v ./... && go install && /bin/bash

EXPOSE 1883
EXPOSE 8300
EXPOSE 8301
EXPOSE 8302
EXPOSE 8400
EXPOSE 8500
EXPOSE 8600