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

Test fixtures are not loaded into development databases. For development, seed files are used. Seed files are preserved and
versioned to maintain healthy development workflows with proper data 
as databases change. 

### How it works

We map seed files to specific migration versions. Seed file names are prefixed with the migration version they map to.

For example, these migration files:
```bash
# Migration files

000001_add_employees_table.down.sql
000001_add_employees_table.up.sql
```
Map to:
```
000001_seed.sql
```
Seed files are placed under `res/seed`.

### Not all migrations will need a seed file

You may have more migrations than seed files because not all migrations require seed data.

For example,
```bash
1_migration 
2_migration
3_migration
4_migration

1_seed 
4_seed
```

#### Rule of thumb

1. If your current migration maps to a seed file, apply it.


2. If your current migration doesn't map to a seed file, __and migrating up is not an option__, migrate down to one that does, 
apply the seed, then migrate back up again.


Continuing the example above, if you want to seed, and your migration version was __3_migration__, without migrating up, you could do:

```bash
 make down 2
 make seed ./res/seed/000001_seed.sql
 make up 2
```

That way you execute a seed file that is guaranteed to work.

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
