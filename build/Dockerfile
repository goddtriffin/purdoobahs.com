FROM golang:1.16.2-alpine3.13 AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /purdoobahs.com
COPY . .
RUN mkdir -p bin
RUN rm -rf bin/*
RUN cp -R ui bin/ui
RUN GOOS=linux GOARCH=amd64 go build -o bin/website ./...

FROM alpine:3.13.2

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /purdoobahs.com
COPY --from=builder /purdoobahs.com/bin/ .

EXPOSE 80

CMD /purdoobahs.com/website -env=production
