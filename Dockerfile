FROM golang:1.13.8-alpine3.11 as builder

WORKDIR /app

# A place to install system dependencies
# e.g. apk add --update --no-cache zsh

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app



FROM scratch

COPY --from=builder /app/app /app/

EXPOSE 8080

ENTRYPOINT ["/app/app"]
