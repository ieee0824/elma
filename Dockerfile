FROM golang:1.9.1 AS build

WORKDIR /go/src/github.com/ieee0824/elma

COPY . .

RUN set -e \
	&& cd cmd/elma \
	&& CGO_ENABLED=0 go build 

FROM alpine:latest

COPY --from=build /go/src/github.com/ieee0824/elma/cmd/elma/elma /bin/elma

CMD ["elma", "-f", "/tmp/config.json"]
