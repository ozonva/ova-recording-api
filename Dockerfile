FROM golang:1.17-alpine  AS builder

RUN apk add --update make

WORKDIR /go/src/github.com/ozonva/ova-recording-api/

COPY . /go/src/github.com/ozonva/ova-recording-api/

RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
RUN mkdir ./config
COPY ./config/config.yml ./config/config.yml
COPY --from=builder /go/src/github.com/ozonva/ova-recording-api/build/ova-recording-api .

RUN chown root:root ova-recording-api

EXPOSE 8081
EXPOSE 8888
CMD ["./ova-recording-api", "--config", "/root/config/config.yml"]
