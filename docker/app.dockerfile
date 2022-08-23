FROM golang:alpine

WORKDIR /pulse-api

ADD . .

EXPOSE 9091 9091

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build "go build -mod=vendor ./cmd/pulse/" -command="./pulse run-api" -polling