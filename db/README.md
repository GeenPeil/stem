
## db
This directory contains files for the database powering this project.
We're using [PostgreSQL](https://www.postgresql.org/)

### model
The database model (`gpstem.dbm`) was created with [pgmodeler](http://www.pgmodeler.com.br/). Please communicate any changes in advance with @GeertJohan to avoid parallel edits. This file cannot be automatically merged with git.

### Setup database
Execute `gpstem.sql` on your postgresql database named "gpstem" to create the database, insert test data, and create a database user:

    # For Ubuntu users:
    $ sudo -u postgres psql
    postgres=# CREATE DATABASE gpstem;
    CREATE DATABASE
    postgres-# \q
    $ sudo -u postgres psql gpstem < gpstem.sql

**Please note: this will create a database user with superuser powers on your database instance, make sure you are not running a public accessible Postgresql database**

### Cleanup

To start from scratch, execute the following statements:

```
REVOKE ALL ON DATABASE gpstem FROM rutte;
DROP USER rutte;
DROP DATABASE gpstem;
```