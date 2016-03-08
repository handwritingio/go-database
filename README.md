# go-database

Provides a wrapper around the go postgres database driver to add migration
management, connection persistence, and transaction composition.

This was extracted from `handwritingio/api`.


## Configuration

The library expects two environment variables to be set:

- `DATABASE_URL`: used for connections created with `GetDB()`.
- `TEST_DATABASE_URL`: used for connections created with `GetTestDB()`.


## Migrations

One folder in your project should contain only raw SQL scripts.

`bootstrap.sql` is a special file that is loaded once onto a clean
database. It sets up the `meta` schema for tracking migrations that have
been applied.

Other files in this directory will be applied to the database in
order determined by their filename. The convention is to name the
file with a version prefix that is zero-padded to 3 digits long
(e.g. `015_index_pizza_size.sql`)

### Hint

While developing a new migration, it's often easiest to guard your
migration file with a `BEGIN` statement at the start of the file and a
`ROLLBACK` at the end. This will let you ensure your syntax is OK and
will throw the migration away at the end. Once you have it working with
no errors, just remove the `BEGIN` and `ROLLBACK` lines.

You can also skip updating the migration log by just piping these scripts
directly to your DB via psql (e.g. `psql mydb < sql/012_drop_foo.sql`)

### Notes

Skipping migration numbers or defining the same version more than once
are both treated as fatal errors to the migration script.


## Code Status

[![Circle CI](https://circleci.com/gh/handwritingio/go-database.svg?style=svg)](https://circleci.com/gh/handwritingio/go-database)
