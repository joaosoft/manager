FROM golang:1.8-alpine

# install git and mercurial
RUN apk add --update git mercurial && rm -rf /var/cache/apk/*
# install make
RUN apk add --update bash make && rm -rf /var/cache/apk/*
# install dep
RUN go get github.com/golang/dep && go install github.com/golang/dep

# install dependencies
ADD Gopkg.toml Gopkg.lock /go/src/go-manager/
RUN cd /go/src/go-manager && dep ensure

# copy configuration
ADD ./config /etc/config

# add source code
ADD . /go/src/go-manager/
WORKDIR /go/src/go-manager/

# copy and build the project
# ADD . /go/src/go-manager/
# RUN cd /go/src/go-manager && dep ensure
# WORKDIR /go/src/go-manager/

EXPOSE 8080
ENTRYPOINT ["go"]
