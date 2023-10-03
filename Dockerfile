FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/ordinals
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/Ordinals /go/src/ordinals


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/Ordinals /usr/local/bin/Ordinals
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["Ordinals"]
