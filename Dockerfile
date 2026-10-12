#
# Trivial Dockerfile .
#
# Build it like so:
#
#   docker build -t http-post-dump:1 .
#

# STEP1 - Build-image
###########################################################################
FROM golang:alpine AS builder

ARG VERSION

LABEL org.opencontainers.image.source=https://github.com/skx/http-post-dump/

# Create a working-directory
WORKDIR $GOPATH/src/github.com/skx/http-post-dump/

# Copy the source to it
COPY . .

# Build the binary - ensuring we pass the build-argument
RUN go build  -o /go/bin/http-post-dump

# STEP2 - Deploy-image
###########################################################################
FROM alpine

# Copy the binary.
COPY --from=builder /go/bin/http-post-dump /usr/local/bin/

# Set entrypoint
ENTRYPOINT [ "/usr/local/bin/http-post-dump" ]

# Create a group and user
RUN addgroup app && adduser -D -G app -h /app app

# Tell docker that all future commands should run as the app user
USER app

# Set working directory
WORKDIR /app
