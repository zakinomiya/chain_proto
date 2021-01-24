FROM golang:1.14-alpine 

RUN apk update && apk add --no-cache alpine-sdk build-base

ENV PJDIR="/go/src/chain_proto"
WORKDIR ${PJDIR}

COPY . .
RUN make server 

CMD ["bin/server", "run"]
