ARG GO_VERSION=1.16
FROM golang:${GO_VERSION} as builder

ENV APP_NAME="main"

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* .
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN go build -mod=readonly -v -o server ./cmd/node/main.go

# Use scratch for clean container.
FROM scratch

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app /

# Run the web service on container startup.
ENTRYPOINT ["/server"]