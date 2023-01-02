FROM alpine:latest
ADD app gateways_config.yaml /app/
RUN chmod +x /app/app
ENTRYPOINT ["/app/app"]