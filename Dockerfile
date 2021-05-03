FROM golang:1.16.3 AS build

WORKDIR /src/
ADD . /src/
RUN make linux

# Use alpine - update SSL certs.
FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Copy the binary.
COPY --from=build /src/bin/aiven /bin/aiven
WORKDIR /bin
ADD websites.csv /bin/websites.csv

ENTRYPOINT ["/bin/aiven"]