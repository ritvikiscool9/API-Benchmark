
FROM golang:1.26 AS builder

WORKDIR /app

# Copy the dependency file first 
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -o api-benchmark: names the output file
# CGO_ENABLED=0: ensures it doesn't need external C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -o api-benchmark .

FROM alpine:latest

# Install certificates so HTTPS requests work
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled binary from the first stage
COPY --from=builder /app/api-benchmark .

ENTRYPOINT ["./api-benchmark"]