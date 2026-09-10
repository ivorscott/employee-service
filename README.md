# Employee Service

Test Go Service

> **Observability lab:** see [`LAB.md`](LAB.md) for a self-paced lab on correlating
> Prometheus metrics to Loki logs in Grafana (data links and correlations).

## Getting started

You need only [docker](https://docs.docker.com/desktop/) and
[make](https://formulae.brew.sh/formula/make):

```bash
make
```

That is the whole setup. It creates `.env` from `.env.sample` if you don't have one,
builds the image, starts every service, runs the database migrations and seeds the
data, then prints where things are listening:

| | |
|---|---|
| API | http://localhost:8080/employees |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| RabbitMQ | http://localhost:15672 (guest/guest) |

```bash
make help    # every available command
make ps      # service status
make logs    # follow logs (make logs S=employee-service for one service)
make stop    # stop, keep data
make down    # remove containers, keep data
make reset   # wipe data and start clean
```

Migrations run automatically when the service boots, and the `seed` job waits for the
service to become healthy before loading `res/seed/data.sql`, so ordering is handled
for you. To reload the seed data on demand: `make seed`.

## Database

```bash
make db               # interactive pgcli shell
make migrate-version  # current migration version
make migrate-up       # apply pending migrations (or N with N=2)
make migrate-down     # roll back 1 migration (or N with N=2)
make migrate-reset    # roll back every migration
make migrate-force N=1        # force a dirty schema to a version, https://bit.ly/3exuENS
make migrate-create NAME=add_widgets   # create a new migration
```

Every one of these runs inside a container on the compose network and talks to postgres
as `employee:5432`, so no `migrate`, `psql` or `pgcli` needs to be installed on your
host. This is deliberate: if you also run postgres natively, it shadows the container on
`localhost:5432` and host-installed tools silently operate on the **wrong** database.

golang-migrate is used internally on service start, in repository tests, and by the
commands above.

## Go development

Working on the Go code (rather than just running the stack) needs a host toolchain:

```bash
make tools     # installs CompileDaemon and mockery
brew install golangci-lint swagger-codegen
```

```bash
make develop   # run the service on the host with hot reload
make test      # unit tests with coverage (starts the test database)
make check     # fmt, vet, lint and test
make build     # compile the binary
```

`make develop` runs the service on your host against the containerised dependencies.
The rest of the stack always runs in containers.

## Troubleshooting

**`password authentication failed`** — postgres only applies `POSTGRES_USER` and
`POSTGRES_PASSWORD` the first time its data volume is created. If you changed them in
`.env` afterwards, the old credentials are still in the volume. Run `make reset`.

**Port already in use** — something else is on 3000, 8080, 9090, 3100, 5432, 5672,
15672 or 12345. Stop it, or change the published port in `docker-compose.yml`.

**No data / 404 from the API** — check the seed job ran: `docker compose logs seed`
should show `INSERT 0 2`. Re-run it with `make seed`.
