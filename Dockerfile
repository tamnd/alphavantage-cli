FROM alpine:3.21

ARG TARGETPLATFORM

RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -H -u 10001 alphavantage

COPY $TARGETPLATFORM/alphavantage /usr/bin/alphavantage

USER alphavantage

ENTRYPOINT ["/usr/bin/alphavantage"]
