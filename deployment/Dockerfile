FROM golang:1.17.7-alpine3.15 AS binary_builder

# update alpine linux dependencies
RUN apk update
RUN apk add --no-cache git ca-certificates tzdata make
RUN update-ca-certificates

WORKDIR /purdoobahs-website

# copy required files
COPY Makefile .
COPY cmd/ cmd/
COPY internal/ internal/
COPY go.mod .
COPY go.sum .

# generate binary
RUN make build

FROM denoland/deno:alpine-1.19.2 as js_builder

# update alpine linux dependencies
RUN apk update
RUN apk add --no-cache make

WORKDIR /purdoobahs-website

# copy required files
COPY Makefile .
COPY ui/static/script/ ui/static/script/
COPY deno.jsonc .

# generate Javascript
RUN make gen_js

FROM node:17.5.0-alpine3.15 as css_builder

# update alpine linux dependencies
RUN apk update
RUN apk add --no-cache make

# install Sass
RUN npm install -g sass

WORKDIR /purdoobahs-website

# copy required files
COPY Makefile .
COPY ui/static/scss/ ui/static/scss/

# generate static assets
RUN make gen_css

FROM alpine:3.15.0

# copy certs
COPY --from=binary_builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=binary_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=binary_builder /etc/passwd /etc/passwd
COPY --from=binary_builder /etc/group /etc/group

WORKDIR /purdoobahs-website

# copy binary
COPY --from=binary_builder /purdoobahs-website/bin/website .

# copy scripts
COPY --from=js_builder /purdoobahs-website/bin/static/script/ static/script/

# copy stylesheets
COPY --from=css_builder /purdoobahs-website/bin/static/stylesheet/ static/stylesheet/

# copy non-generative static assets
COPY assets/ assets/
COPY ui/html/ html/
COPY ui/static/file/ static/file/
COPY ui/static/image/ static/image/
COPY ui/static/video/ static/video/

# run server
EXPOSE 8080
CMD ./website -env=production
