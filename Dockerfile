FROM golang:1.16.3 AS build

WORKDIR /src/
ADD . /src/
RUN make linux

# Use alpine - update SSL certs.
FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Copy the binary.
RUN mkdir /app
COPY --from=build /src/bin/aiven /app/aiven
WORKDIR /app
ADD websites.csv /app/websites.csv

ENTRYPOINT ["/app/aiven"]