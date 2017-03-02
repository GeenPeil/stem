#!/bin/bash
set -e

echo "Drop existing database"
sudo -i -u postgres psql postgres <<EOF
REVOKE ALL ON DATABASE gpstem FROM rutte;
DROP USER IF EXISTS rutte;
DROP DATABASE IF EXISTS gpstem;
CREATE DATABASE gpstem;
EOF

echo "Creating database structure"
sudo -i -u postgres psql -v ON_ERROR_STOP=1 gpstem < gpstem.sql | sed 's/^/[schema gpstem.sql] /'

echo "Inserting i8n data"
datafiles=(\
	01-i8n-countries\
	02-i8n-languages\
	03-i8n-country_names-nl_NL\
	04-i8n-country_names-en_US\
)
for d in ${datafiles[@]}; do
	sudo -i -u postgres psql -v ON_ERROR_STOP=1 gpstem < data/${d}.sql | sed 's/^/[data '"${d}"'.sql] /'
done
