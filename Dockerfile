FROM golang:1.24-alpine

WORKDIR /app
COPY ./ ./
RUN go mod download

WORKDIR /app/cmd/api-server
# CMD [ "go","run","main.go" ]