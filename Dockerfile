# Backend service image
FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w" -o /out/blog-server ./cmd/server

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && mkdir -p /app/configs /app/scripts /app/logs /app/uploads

ENV TZ=Asia/Shanghai

COPY --from=builder /out/blog-server /app/blog-server
COPY configs/config.example.yaml /app/configs/config.yaml
COPY scripts/ /app/scripts/

EXPOSE 9090

CMD ["/app/blog-server"]
