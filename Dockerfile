FROM golang:1.22.4-alpine

RUN apk update && apk add git

ENV CGO_ENABLED=0
ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/zeus

COPY . .

RUN go mod download
RUN GOOS=linux go build -o app
ENTRYPOINT ["./app"]

EXPOSE 1994