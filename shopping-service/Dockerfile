FROM golang:1.23-alpine

WORKDIR /app
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

COPY . .
RUN go build -o shopping .

CMD ["./shopping"]
