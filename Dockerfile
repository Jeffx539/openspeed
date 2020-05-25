FROM node:buster

RUN apt-get update && apt-get install -y golang

ENV GOPATH=/go

RUN mkdir -p /go/src/github.com/jeffx539/openspeed
WORKDIR /go/src/github.com/jeffx539/openspeed
COPY . /go/src/github.com/jeffx539/openspeed

RUN go get && cd web yarn install && cd .. && make &&  mkdir -p /openspeed && mv bin/* /openspeed && curl https://iptoasn.com/data/ip2asn-combined.tsv.gz -o /openspeed/ip2asn-combined.tsv.gz

WORKDIR /openspeed
ENTRYPOINT ["/openspeed/openspeed"]
