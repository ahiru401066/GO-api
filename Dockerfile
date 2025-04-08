FROM golang:latest

COPY ./ ./

CMD [ "bash" ]
# CMD [ "go","run","main.go" ]