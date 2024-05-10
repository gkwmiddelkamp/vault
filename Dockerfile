FROM golang:1.22-alpine as build

WORKDIR /src
COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /vault
RUN addgroup vault &&  \
     adduser --ingroup vault --uid 1000 --shell /bin/false --disabled-password vault && \
     cat /etc/passwd | grep vault > /etc/passwd_vault && \
     cat /etc/group | grep vault > /etc/group_vault

FROM scratch as vault
WORKDIR /
COPY --from=build /vault /usr/local/bin/vault
COPY --from=build /etc/passwd_vault /etc/passwd
COPY --from=build /etc/group_vault /etc/group
USER vault:vault
EXPOSE 8080
CMD ["vault"]