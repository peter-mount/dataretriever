# Dockerfile used to build the application
FROM golang:latest

ADD app /app/
WORKDIR /app
RUN go build -o bridge .

CMD ["/app/bridge"]
