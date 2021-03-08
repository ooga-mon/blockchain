ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* .
RUN go mod download

# Copy local code to the container image.
COPY . .

#disable crosscompiling 
ENV CGO_ENABLED=0

#compile linux only
ENV GOOS=linux

# Build the binary.
RUN go build  -ldflags '-w -s' -a -installsuffix cgo -o server ./cmd/node/main.go

# Use slim for smaller container size.
FROM scratch

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app /

# Run the web service on container startup.
ENTRYPOINT ["/server"]