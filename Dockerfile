FROM golang:1.14-alpine AS builder
# Create appuser
# ENV USER=appuser
# ENV UID=10001
ENV PJDIR="/go/src/chain_proto"
RUN apk update && apk add --no-cache make gcc
RUN apk add --no-cache alpine-sdk build-base

COPY . ${PJDIR}
RUN cd ${PJDIR} && make server 
#RUN cd ${PJDIR} && touch bin/server 

#FROM alpine
#ENV PJDIR="/go/src/chain_proto"
#COPY --from=builder ${PJDIR}/bin/server /go/bin/server

#CMD ["/go/bin/server", "run"]

CMD ["/go/src/chain_proto/bin/server", "run"]
