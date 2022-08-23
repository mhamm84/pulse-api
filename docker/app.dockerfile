FROM golang:alpine
RUN apk add --no-cache ca-certificates git
RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /pulse-api

ADD . .

EXPOSE 9091 9091

#RUN go install -mod=mod github.com/githubnemo/CompileDaemon
#RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build "go build -mod=vendor ./cmd/pulse/" -command="./pulse run-api" -polling