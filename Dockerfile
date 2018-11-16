FROM golang:1.11

ARG BRCH

ENV BRAN $BRCH

WORKDIR /go/src/github.com/aeternas/SwadeshNess-IntegrationTests
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM frolvlad/alpine-bash:latest

ARG BRCH

ENV BRANCH $BRCH

RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/aeternas/SwadeshNess-IntegrationTests .

CMD ["bash"]
