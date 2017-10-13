FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD trumail /usr/local/bin/trumail
ADD /web /web/
EXPOSE 8000
CMD ["trumail"]