FROM golang:1.16.3 AS build

WORKDIR /src/
ADD . /src/
RUN make linux

FROM scratch
COPY --from=build /src/bin/aiven /bin/aiven

ENTRYPOINT ["/bin/aiven"]