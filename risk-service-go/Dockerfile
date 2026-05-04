#uses official go image to build and run the service
FROM golang:1.24.1

#Sets the working directory inside the container
WORKDIR /app

#copies the Go module file first
COPY go.mod ./

#copies the rest of the source code
COPY . .

#builds the go app into a binary named risk-service
RUN go build -o risk-service main.go

#documents that the container listens on port 8081
EXPOSE 8081

#starts the compiled service when the container runs
CMD ["./risk-service"]