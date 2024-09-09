FROM --platform=$BUILDPLATFORM golang:1.23.0 AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./tests tests

# Run tests
RUN go test ./... -v

# Set environment variables for the target platform
ARG TARGETOS
ARG TARGETARCH

# Build the Go application
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o booking_server ./cmd/server/main.go

# Use a minimal base image for the final image
FROM --platform=$TARGETPLATFORM alpine:3.20.0

# Install required packages including tzdata
RUN apk --no-cache add ca-certificates tzdata

# Set the timezone
ENV TZ=Asia/Bangkok
RUN ln -sf /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && echo "Asia/Bangkok" > /etc/timezone

# Set the working directory
WORKDIR /root/

# Copy the Go binary from the build stage
COPY ./configs ./configs
COPY --from=build /app/booking_server .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./booking_server"]
