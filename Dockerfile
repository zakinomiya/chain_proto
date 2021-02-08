FROM golang:1.14-alpine AS builder

RUN apk update && apk add --no-cache alpine-sdk build-base

ENV PJDIR="/go/src/chain_proto"
WORKDIR ${PJDIR}

COPY . .
RUN make server 

#FROM scratch
FROM alpine

ENV GOPATH="/go"
WORKDIR ${GOPATH}/src/chain_proto
COPY --from=builder /go/src/chain_proto/bin ./bin
COPY --from=builder /go/src/chain_proto/db/sql ./db/sql
COPY --from=builder /go/src/chain_proto/config ./config

CMD ["bin/server", "run"]
