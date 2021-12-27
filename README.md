# Employee Service

__Required__
- docker
- golangci-lint
- CompileDaemon
- swagger-codegen

Test Go Service

## Usage

```bash
make # start app
docker-compose up -d # start containers
```

The employee service does not run in a container during local development.
Docker containers are only leveraged for database and observability services.

# To Do 

[X] add logging w/ elk stack

[X] log retention policy

[X] add swagger

[X] add zappr

[X] add golangci-lint

[X] add monitoring w/ prometheus and grafana

[X] add tracing w/ Open telemetry

[ ] add postgres

[ ] add crud 

[ ] add interfaces

[ ] add fixtures

[ ] add tests

[ ] add goldenfiles

[ ] add test data

[ ] auth

[ ] cors

[ ] github workflow

[ ] terraform 

[ ] kubernetes