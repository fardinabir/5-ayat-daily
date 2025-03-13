FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./main

FROM alpine AS app
WORKDIR /app
COPY --from=builder /app/main .
COPY config/example-config.yaml ./config/.config.yaml
EXPOSE 8085
CMD ["./main"]