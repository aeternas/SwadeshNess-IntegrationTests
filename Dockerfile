FROM golang:1.10

ARG YANDEX_API_KEY=foo

ENV YANDEX_API_KEY=${YANDEX_API_KEY}

WORKDIR /go/src/github.com/aeternas/SwadeshNess-IntegrationTests
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/aeternas/SwadeshNess-IntegrationTests .

EXPOSE 8080

CMD ["./SwadeshNess-IntegrationTests"]
