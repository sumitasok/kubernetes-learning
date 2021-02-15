FROM golang:1.15.8-buster as buildstage

WORKDIR /go/src/github.com/sumitasok/kubernetes-learning

ADD . .

WORKDIR /go/src/github.com/sumitasok/kubernetes-learning/cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build -o server

#---------#---------#---------#---------#---------#---------#---------#

FROM alpine:3.5

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --no-cache tzdata

ENV GOPATH /go

WORKDIR /go/src/github.com/sumitasok/kubernetes-learning

COPY --from=buildstage /go/src/github.com/sumitasok/kubernetes-learning/cmd/server .

RUN mkdir /store

CMD ["./server"]