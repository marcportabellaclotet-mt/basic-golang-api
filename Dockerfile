FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 -S api && adduser -u 1001 -S api  -G api
RUN pass=$(echo date +%s | sha256sum | base64 | head -c 32; echo | mkpasswd) && \
    echo "api:${pass}" | chpasswd

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group
COPY --from=alpine /etc/shadow /etc/shadow

COPY --chown=api build/api_linux_amd64 /api

USER api
CMD ["/api"]