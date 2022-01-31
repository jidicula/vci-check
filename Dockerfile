FROM golang:1.16-alpine as builder

RUN apk --no-cache add ca-certificates
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/vci-check .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/vci-check /bin/vci-check
ENTRYPOINT [ "/bin/vci-check" ]
