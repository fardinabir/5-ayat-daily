FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./main

FROM golang:1.19-alpine AS app
WORKDIR /app
COPY --from=builder /app/main .
COPY config/example-config.yaml ./config/.config.yaml
EXPOSE 8085
CMD ["./main"]