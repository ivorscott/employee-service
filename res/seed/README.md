# Seed Versioning


Data seeding is the process of populating a database with an initial set of data.
This folder contains seed files. __Seeding is for development__ environments.

## Naming convention

Databases change. We leave seeding out of migrations because migration files
will be used in production environments. In order to prevent seed files from becoming a burden
as databases change, we associate new seed files to specific migration versions.

For example, seed files prefixed `000001` indicate they are designed to work with migration `version 1`:

```bash
# Seed files

000001_seed_employees_table.down.sql
000001_seed_employees_table.up.sql
```

```bash
# Migration files

000001_add_employees_table.down.sql
000001_add_employees_table.up.sql
```

Before applying a seed file, ensure you are on the correct migration 
version to guarantee it works. 

``` 
make version
```

Apply the seed.

```
make seed ./res/dataseed/000001_seed_employees_table.up.sql
```

