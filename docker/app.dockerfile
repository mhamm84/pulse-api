###################################################
# Plain env as basement and for local development #
###################################################
FROM golang:alpine as env

RUN apk add --no-cache ca-certificates git

FROM env as dev
# Hot reload using CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /pulse-api
ADD . .

EXPOSE 9091 9091
ENTRYPOINT CompileDaemon -build "go build -mod=vendor ./cmd/pulse/" -command="./pulse run-api" -polling

##########################################################
# Prepare a build container with all dependencies inside #
##########################################################
FROM env as builder

WORKDIR /pulse-api
ADD . .

RUN go build -o /go/bin/pulse -mod=vendor ./cmd/pulse/

###########################################
# Create clean container with binary only #
###########################################
FROM alpine as exec

RUN apk add --update bash ca-certificates

WORKDIR /app
COPY --from=builder /go/bin/pulse ./

EXPOSE 9091 9091
CMD ["./pulse", "run-api"]