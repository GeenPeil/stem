
## db
This directory contains files for the database powering this project.
We're using [PostgreSQL](https://www.postgresql.org/)

### model
The database model (`gpstem.dbm`) was created with [pgmodeler](http://www.pgmodeler.com.br/). Please communicate any changes in advance with @GeertJohan to avoid parallel edits. This file cannot be automatically merged with git.

### Setup database
Execute `gpstem.sql` on your postgresql database to create the database, insert test data, and create a database user.
**Please note: this will create a database user with superuser powers on your database instance, make sure you are not running a public accessible Postgresql database **