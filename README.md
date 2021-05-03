aiven
============

Small test project.

What do you need?

1. Docker
2. `make`
3. `git`
4. `golang` - not necessary if you only use docker
5. Kafka Service from Aiven.io.
6. Postgres Service from Aiven.io.

How to launch?

```bash
git clone https://github.com/darron/aiven.git
cd aiven
# Download service.cert, service.key, ca.pem from Aiven Kafka service page.
# Copy them into this folder.
# Update docker.env with KAFKA_HOST and POSTGRES_URL from Aiven services.
cp docker.env.dist docker.env
make docker
docker compose up
```