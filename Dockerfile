# Use a specific version of Alpine
FROM alpine:3.14

WORKDIR /app
COPY retail-pulse .env /app/
RUN chmod +x /app/retail-pulse
CMD ["./retail-pulse"]


