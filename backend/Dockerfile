# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy all source code first
COPY . .

# Run go mod tidy to generate go.sum since we have all source code now
RUN go mod tidy

# Build single binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sn-backend cmd/server/main.go

# Run stage
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app
COPY --from=builder /app/sn-backend .

EXPOSE 8080
CMD ["./sn-backend"]
