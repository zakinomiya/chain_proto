FROM golang:1.14-alpine as builder
# Create appuser
# ENV USER=appuser
# ENV UID=10001
ENV PJDIR=/go/src/chain_proto
RUN apk update && apk add --no-cache make gcc
RUN apk add --no-cache alpine-sdk build-base

COPY . ${PJDIR}
RUN cd ${PJDIR} && make server 

FROM scratch

COPY --from=builder ${PJDIR}/bin/server /go/bin/server

ENTRYPOINT ["/go/bin/server"]

