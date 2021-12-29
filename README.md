# Employee Service

Test Go Service

__Required__
- docker
- golangci-lint
- CompileDaemon
- swagger-codegen
- golang-migrate
- psql

__Optional__
- pgcli

## Usage

```bash
make # start app
docker-compose up -d # start containers
```

## Migration and Seeding

```bash

make migration <name> # create migration

make version # print version

make up # migrate to latest migration

make up <num> # migrate up N migrations

make down  # migrate down 1 migration

make down <num>  # migrate down N migrations

make downfall # apply all down migrations

make force <version> # force a version https://bit.ly/3exuENS

make seed <filepath> # seed database
```
### Seed Versioning
We map seed files to specific migration versions. [Learn more](./res/seed/README.md).
## Entering Database

```bash
make db # enter postgres database 
```

The employee service does not run in a container during local development.
Containers are only used for databases and observability services.

# To Do 

[X] add logging w/ elk stack

[X] log retention policy

[X] add swagger

[X] add zappr

[X] add golangci-lint

[X] add monitoring w/ prometheus and grafana

[X] add tracing w/ Open telemetry

[X] add postgres

[X] add makefile

[X] add interfaces

[ ] add fixtures

[ ] add tests

[ ] add goldenfiles

[ ] add crud

[ ] add test data

[ ] auth

[ ] cors

[ ] github workflow

[ ] terraform 

[ ] kubernetes