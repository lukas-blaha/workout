#!/bin/bash

# Create table if not exists
echo "Creating table..."
docker exec -ti workout-postgres-1 psql -U postgres -d exercises -c "create table if not exists exercises(id serial primary key, name varchar(50), count int, date varchar(10))" >/dev/null 2>&1

sleep 5

# Populate DB with example data
echo "Sending example data to db..."
curl http://localhost:8082/exercises -X POST -d @example.json >/dev/null 2>&1 && \
echo "Done"
