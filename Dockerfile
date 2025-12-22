FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ezkitchen ./cmd/web

EXPOSE 4000

CMD ["./ezkitchen"]
