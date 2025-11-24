FROM golang:1.24.5-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o product-service
CMD ["./product-service"]