FROM golang:1.21.4

WORKDIR /app

COPY ./ ./
RUN go mod download


RUN go build -o wiskundi-app ./cmd/web


EXPOSE 8080


CMD ["./wiskundi-app"]