FROM golang:1.9.1 AS build

WORKDIR /go/src/github.com/ieee0824/elma

COPY . .

RUN set -e \
	&& cd cmd/elma \
	&& CGO_ENABLED=0 go build 

FROM alpine:latest

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=build /go/src/github.com/ieee0824/elma/cmd/elma/elma /bin/elma

CMD ["elma", "-f", "/tmp/config.json"]
