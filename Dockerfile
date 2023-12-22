FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN go build -o /bin/NSCMBot ./cmd/main.go

FROM scratch
COPY data/ data/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/NSCMBot ./NSCMBot
EXPOSE ${WEBHOOK_PORT}
ENTRYPOINT ["./NSCMBot"]

