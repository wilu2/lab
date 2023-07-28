FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    TZ=Asia/Shanghai

WORKDIR /app
