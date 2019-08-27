FROM golang:1.12-alpine AS build

ARG which=backendA

COPY . /src
WORKDIR /src/$which

RUN apk add --no-cache git build-base \
    && go build -o /backend

FROM alpine

COPY --from=build /backend /backend

USER nobody

EXPOSE 8080/tcp

CMD ["/backend"]
