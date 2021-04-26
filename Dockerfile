FROM golang:latest as builder
RUN mkdir -p /go/src/faker
WORKDIR /go/src/faker
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/faker/app .
COPY --from=builder /go/src/faker/config.json .
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&
    touch cron.daily &&
    echo "1 1,3 * * * ./app" >> cron.daily &&
    cat cron.daily >> /var/spool/cron/crontabs/root
CMD ["crond"]
