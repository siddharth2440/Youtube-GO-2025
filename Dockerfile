FROM golang:latest
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

EXPOSE 7565
CMD [ "air", "-c",".air.toml" ]