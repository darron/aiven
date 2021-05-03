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
# Copy them into this folder.
cp docker.env.dist docker.env
# Update docker.env with KAFKA_HOST and POSTGRES_URL from Aiven services.
make docker
docker compose up
```

What can be better?

- [ ] Dependency injection to help with mocking
- [ ] Additional retries and error checking for Kafka writes
- [ ] Additional retries and error checking for Postgres writes
- [ ] More tests of all varieties: unit, integration, mocking
- [ ] Try out https://github.com/testcontainers/testcontainers-go for end to end tests
- [ ] Naive datbase schema improved
- [ ] Move to protobuf for Kafka transport to optimize
- [ ] Expose metrics for the running services via HTTP
- [ ] Putting certificates inside Docker images isn't great - inject at runtime