FROM golang:1.23
LABEL maintainer="henny@krijnen.dev"
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /google-logging-hmac-proxy

EXPOSE 8080

ENTRYPOINT ["/google-logging-hmac-proxy"]