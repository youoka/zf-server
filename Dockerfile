FROM golang:1.24.0 AS builder
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go mod tidy
WORKDIR /app/cmd
RUN go build -o main .

FROM alpine:latest
ENV TZ=Asia/Shanghai
WORKDIR /root/
COPY --from=builder /app/cmd/main .
COPY --from=builder /app/config/config.yaml .
COPY --from=builder /app/static/. ./static
CMD ["./main"]