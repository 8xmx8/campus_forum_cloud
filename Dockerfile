FROM golang:1.22.5-alpine AS builder

LABEL stage=gobuilder
LABEL authors="yangjianyang"

ENV CGO_ENABLED 0
# http://172.16.20.30:32888 为自建的go第三方库代理
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /app

COPY . .
# RUN git clone -b v0.1.4 --depth=1 ssh://git@git.gobies.org:22703/sutra/gosutra.git scripts/gosutra
# RUN git clone -b v0.0.3 --depth=1 ssh://git@git.gobies.org:22703/fofa-backend/aho-corasick.git scripts/aho-corasick
RUN go mod download
# 交叉编译
RUN go build -ldflags="-s -w" -o campus main.go

FROM alpine:3.17
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /app/config /app/config

ENTRYPOINT ["./campus"]
 CMD ["-f", "config/config.conf"]