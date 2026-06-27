FROM golang:1.25-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend-go .

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build /app/backend-go /app/backend-go

EXPOSE 8080
CMD ["/app/backend-go"]
