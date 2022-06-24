# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . ./

RUN go build -ldflags -s -v -o das-indexer cmd/main.go

##
## Deploy
##
FROM ubuntu

ARG TZ=Asia/Shanghai

RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get install -y ca-certificates tzdata \
    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/das-indexer /app/das-indexer
COPY --from=build /app/config/config.yaml /app/config/config.yaml

EXPOSE 8121 8122 8123

ENTRYPOINT ["/app/das-indexer", "--config", "/app/config/config.yaml"]
