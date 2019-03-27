FROM golang:alpine as builder
RUN apk add git
RUN mkdir -p $GOPATH/src/github.com/simelo
RUN cd $GOPATH/src/github.com/simelo && git clone https://github.com/simelo/rextporter.git
RUN cd $GOPATH/src/github.com/simelo/rextporter/cmd/rextporter && go install ./...


FROM alpine:latest as builder
WORKDIR /bin
COPY --from=build_layer /go/bin/rextporter .