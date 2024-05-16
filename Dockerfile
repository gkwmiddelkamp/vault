FROM golang:1.22-alpine as build

WORKDIR /src
COPY . .
RUN apk update &&  \
    apk upgrade &&  \
    apk add --no-cache ca-certificates && \
    update-ca-certificates
RUN go mod download &&  \
    go mod verify &&  \
    go build -v -o /vault

FROM scratch as vault
WORKDIR /
USER 10001:10002
EXPOSE 8080
COPY --from=build /vault /usr/local/bin/vault
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["vault"]