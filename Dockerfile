FROM registry.intsig.net/acg-sre/golang:1.18.5-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./configs/api-server.yaml /app/configs/api-server.yaml
COPY ./configs/system-token /app/configs/system-token
RUN go build -ldflags="-s -w" -o /app/textin cmd/api-server/apiserver.go


FROM registry.intsig.net/acg-sre/alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/textin /app/textin
COPY --from=builder /app/configs /app/configs

CMD ["./textin", "-c", "configs/api-server.yaml"]
