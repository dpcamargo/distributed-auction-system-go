FROM golang:1.23-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-arm64} go build -ldflags="-w -s" -o auction cmd/auction/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auction /app/auction
COPY cmd/auction/.env /app/.env

EXPOSE 8080

ENTRYPOINT ["/app/auction"]
