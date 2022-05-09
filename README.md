# Employee Service

Test Go Service

__Required__
- [make](https://formulae.brew.sh/formula/make)
- [docker](https://docs.docker.com/desktop/)
- [golangci-lint](https://formulae.brew.sh/formula/golangci-lint)
- [CompileDaemon](https://github.com/githubnemo/CompileDaemon)
- [swagger-codegen](https://formulae.brew.sh/formula/swagger-codegen)
- [golang-migrate](https://formulae.brew.sh/formula/golang-migrate)
- [psql](https://formulae.brew.sh/formula/postgresql)
- [mockery](https://github.com/vektra/mockery)

__Optional__
- [pgcli](https://formulae.brew.sh/formula/pgcli)


## Usage


Clone `.env.sample` and rename it `.env`.

```bash
docker-compose up -d # start containers
make # start app
make test
```

## Migration and Seeding

The service won't have any data until you seed its database.

```bash
make seed
```

### Migration Commands

golang-migrate is used internally upon service start, in repository tests, and through makefile helper commands. 
```bash

make migration <name> # create migration

make version # print version

make up # migrate to latest migration

make up <num> # migrate up N migrations

make down  # migrate down 1 migration

make down <num>  # migrate down N migrations

make downfall # apply all down migrations

make force <version> # force a version https://bit.ly/3exuENS

```

## Entering Database

```bash
make db # enter postgres database 
```

The service does not run in a container during local development.
Containers are only used for databases and observability services.

