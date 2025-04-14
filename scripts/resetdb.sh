#!/bin/bash
psql -U postgres -p 5433 -c "DROP DATABASE IF EXISTS 2fa"
psql -U postgres -p 5433 -c "CREATE DATABASE 2fa"
echo "Database reset complete"