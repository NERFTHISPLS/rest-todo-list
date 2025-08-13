FROM golang:alpine

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
