FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
EXPOSE 8080
CMD ["go", "run", "main.go"]
