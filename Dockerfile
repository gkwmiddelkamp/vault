FROM golang:1.22-alpine as build

WORKDIR /src
COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /vault

FROM scratch as vault
WORKDIR /
USER 1000
EXPOSE 8080
COPY --from=build /vault /usr/local/bin/vault
ENTRYPOINT ["vault"]