FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /app

COPY . .

RUN go build -o app ./main.go

FROM scratch
COPY --from=builder /app .
EXPOSE 8100
ENTRYPOINT ["/app"]