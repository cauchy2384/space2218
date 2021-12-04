FROM golang:1.17.4-alpine3.14 AS builder

WORKDIR /space2218

# Install make
RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache make

COPY . /space2218
RUN make build 

FROM alpine:3.14  
WORKDIR /app
COPY --from=builder /space2218/app /app