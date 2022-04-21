ARG APP_NAME
ARG GOLANG_VERSION

# used to generate the application binary
FROM golang:${GOLANG_VERSION}

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN rm -rf vendor Gopkg.* ; exit 0
RUN dep init ; exit 0
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/${APP_NAME} examples/main.go




# used to create the application image
FROM alpine

WORKDIR /bin/${APP_NAME}

COPY --from=0 go/src/${APP_NAME}/bin/ .

EXPOSE 8080

CMD ["./${APP_NAME}"]