FROM quay.io/projectquay/golang:1.20 as builder


# argument from makefile to set type of arch
ARG TARGETARCH  

WORKDIR /go/src/app
COPY . .

RUN make build TARGETARCH=$TARGETARCH

FROM scratch
WORKDIR /

COPY --from=builder /go/src/app/sbot .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./sbot"]
CMD [ "start" ]