#Define the base image
FROM golang:1.23.5-alpine AS builder

#set the working directory
WORKDIR /app

#copy go.mod and go.sum file
COPY ./ ./

#downloads go dependencies
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
# copy source files
# COPY ./ .

#build the go app
RUN go build -o main

#use a smaller image to run the app
FROM alpine:latest

#set the working directory
WORKDIR /root/

#copy the compiled go binary from the builder image
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs


EXPOSE 9091

CMD ["./main"]