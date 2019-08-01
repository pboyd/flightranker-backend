FROM golang:1.12-alpine

ARG which=backendA

COPY . /src
WORKDIR /src/$which

RUN apk add --no-cache --virtual .deps git build-base \
    && go build -o backend \
    && mv backend / \
    && rm -rf /go /usr/local/go /src \
    && apk del .deps

EXPOSE 8080/tcp

ENTRYPOINT ["/backend"]
