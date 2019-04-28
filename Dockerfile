FROM alpine as builder

RUN apk update && apk add ca-certificates

FROM scratch

COPY main /http-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
VOLUME /http

CMD ["/http-server", "start", "/http"]
