FROM golang:1.21rc2-alpine3.18
RUN go version
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
# Set the working directory inside the container
WORKDIR /app

# Set the container's entry point to run the built application
