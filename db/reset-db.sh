#!/bin/bash
set -e

sudo -u postgres psql postgres <<EOF
REVOKE ALL ON DATABASE gpstem FROM rutte;
DROP USER rutte;
DROP DATABASE gpstem;
CREATE DATABASE gpstem;
EOF

sudo -u postgres psql gpstem < gpstem.sql