FROM golang:1.19 as builder

ARG CGO_ENABLED=0

WORKDIR /app

COPY . .
RUN go mod download

RUN make compile

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/osoc-dialog /osoc-dialog
COPY --from=builder /app/swagger/api.swagger.json /swagger/api.swagger.json
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip

ENTRYPOINT ["/osoc-dialog"]
