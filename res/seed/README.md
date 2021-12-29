# Seed Versioning

Databases change. For ease of use, we version seed files as our databases change.

## How it works

We map seed files to specific migration versions and seed files don't depend on each other. 

Each seed file MUST:

1) EMPTY all available tables
2) SEED all available tables.
3) MAP to a migration

This approach may create duplicate content overtime, but it removes dependencies. 
Seed file names are prefixed with the migration version they map to. __Some migrations won't need a seed file__.

For example,
```bash
# Migration files

000001_add_employees_table.down.sql
000001_add_employees_table.up.sql
```
The corresponding seed file name would be: `000001_seed.sql`. All seed files MUST be placed under `res/seed/`.

Apply the seed.

```bash
make version # check your version
make seed ./res/seed/000001_seed.sql
```

### Not all migrations will need a seed file

You may have more migrations than seed files because not all migrations require new seed data. 

For example,
```bash
1_migration 
2_migration
3_migration
4_migration

1_seed 
4_seed
```

__If a seed file exists for your current migration apply it.__

