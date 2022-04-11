FROM golang:1.18

WORKDIR /app

# download dependent packages
COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
