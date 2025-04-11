FROM golang:latest

WORKDIR /app
COPY ./ ./
RUN go mod download

CMD [ "bash" ]
# CMD [ "go","run","main.go" ]