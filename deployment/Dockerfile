FROM golang:1.17.7-alpine3.15 AS binary_builder

# update alpine linux dependencies
RUN apk update
RUN apk add --no-cache git ca-certificates tzdata make
RUN update-ca-certificates

WORKDIR /purdoobahs-website

# generate binary
COPY . .
RUN make build

FROM node:17.5.0-alpine3.15 as static_assets_builder

# update alpine linux dependencies
RUN apk update
RUN apk add --no-cache make

# install Sass
RUN npm install -g sass

WORKDIR /purdoobahs-website

# generate binary
COPY . .
RUN make gen_static

FROM alpine:3.15.0

# copy certs
COPY --from=binary_builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=binary_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=binary_builder /etc/passwd /etc/passwd
COPY --from=binary_builder /etc/group /etc/group

WORKDIR /purdoobahs-website

# copy binary
COPY --from=binary_builder /purdoobahs-website/bin/ .

# copy static assets
COPY --from=static_assets_builder /purdoobahs-website/bin/ .

# run server
EXPOSE 8080
CMD ./website -env=production
