# Golang Server's Dockerfile
FROM golang:1.22.1-alpine

# Set the working directory
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main .

CMD ["/app/main"]