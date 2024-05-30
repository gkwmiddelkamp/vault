FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN apk update &&  \
    apk upgrade &&  \
    apk add --no-cache ca-certificates && \
    update-ca-certificates
RUN go mod download &&  \
    go mod verify &&  \
    go mod vendor && \
    go build -v -o /vault
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

FROM scratch AS vault
WORKDIR /
USER 10001:10002
EXPOSE 8080
COPY --from=build /vault /usr/local/bin/vault
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/docs/swagger.json /docs/swagger.json
ENTRYPOINT ["vault"]