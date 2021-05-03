aiven
============

Small test project.

What do you need?

1. Docker
2. `make`
3. `git`
4. `golang` - not necessary if you only use Docker
5. Kafka Service from Aiven.io.
6. Create `websites` topic in Kafka.
6. Postgres Service from Aiven.io.

How to launch?

```bash
git clone https://github.com/darron/aiven.git
cd aiven
# Download service.cert, service.key, ca.pem from Aiven Kafka service page.
# Copy them into the "certs" folder.
cp docker.env.dist docker.env
# Update docker.env with KAFKA_HOST and POSTGRES_URL from Aiven services.
# Build the image locally.
make docker
# NOTE: If you don't want to build the image locally, change: "aiven:latest"
# in docker-compose.yml to "darron/aiven:latest": 
# https://hub.docker.com/repository/docker/darron/aiven
docker compose up --always-recreate-dep # On Linux you might need to run "docker-compose"
```

What can be better?

- [x] Dependency injection to help with mocking
- [x] Additional retries and error checking for Kafka writes
- [x] Additional retries and error checking for Postgres writes
- [x] Putting certificates inside Docker images isn't great - inject at runtime
- [ ] More tests of all varieties: unit, integration, mocking
- [ ] Try out https://github.com/testcontainers/testcontainers-go for end to end tests
- [ ] Naive datbase schema improved - could optimize, do some normalization, add rollups, use table partitioning
- [ ] Move to protobuf for Kafka transport to optimize
- [ ] Expose metrics for the running services via HTTP
- [ ] Use goroutines so that write operations don't block
- [ ] If we get an error while getting metrics - add a blank metric to the DB.