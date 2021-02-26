FROM golang:1.16.0-alpine AS build_base

RUN apk add --no-cache git \
                        build-base \
                        make \
                        upx

ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /tmp/pokemonb2w

# We want to populate the module cache based on the go.{mod,sum} files.
COPY . .
RUN make mod

# Generate swagger static UI files
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
RUN make swagger

# Unit tests
RUN make test

# Build the Go app
RUN make build
RUN make pack

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

# Copy binary
COPY --from=build_base /tmp/pokemonb2w/bin/pokemonb2w /app/pokemonb2w

# Copy swagger static files
COPY --from=build_base /tmp/pokemonb2w/static /static

# Run the binary program produced by `build phase`
CMD /app/pokemonb2w

# This container exposes port 8080 to the outside world
EXPOSE 8080
