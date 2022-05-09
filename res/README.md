# Res

Repository config folder and package

- [Config](#config)
- [Fixtures](#fixtures)
- [Seeds](#seeds)
- [Migrations](#migrations)
- [Golden Files](#golden-files)

Repository concerns include configuration for data services, test fixtures, golden files, migrations and seeding the database.  

```bash
config # configuration for data services
fixtures # test fixtures
golden # golden files
migrations # migrations
seed # data for development
```

## Config

`res/config`

Data services may require configuration yaml. This configuration is preserved here.

## Fixtures

`res/fixtures`

Test fixtures are only loaded into test databases. Fixtures data is feed into our repository level tests. Fixtures allow 
the Go service to be tested against a real database instead of running it against mocks, which may lead to production bugs
not being caught in the tests.

__Before every test, the test database is cleaned and the fixture data is loaded into
the database.__ https://github.com/go-testfixtures/testfixtures

## Seeds

`res/seed`

Seed data should be updated as databases change.

### How it works

Each table should have a corresponding seed file.

## Migrations

`res/migrations`

Migrations are managed via `res.MigrateUp()` and via 
`make` commands. [Learn more](/README.md#migration-and-seeding).

## Golden Files

`res/golden`

Goldenfiles are used in tests to compare database responses with previous queries preserved as snapshots in json format.
If a database response changes, the golden file test fails and a new snapshot must be saved for the test to pass. 
To update all golden files run:
```
go test ./... -update
```
Alternatively, if you want one golden file to update, comment the corresponding
code block:

```go
// pkg/repository/repository_test.go

goldenFile := "employee.json"

//if golden.ShouldUpdate() {
    testutils.SaveGoldenFile(&actual, goldenFile)
//}
```
