FROM docker.1ms.run/library/golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/seckill ./cmd/seckill/

FROM docker.1ms.run/library/alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /build/seckill .
COPY configs/config.json .

EXPOSE 8888

ENTRYPOINT ["./seckill"]
