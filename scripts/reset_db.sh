#!/bin/bash

export CONTAINER_NAME=octovy-db

if [ -z "$POSTGRES_USER" ]; then
    echo "Error: POSTGRES_USER environment variable is not set"
    exit 1
fi

if [ -z "$POSTGRES_PASSWORD" ]; then
    echo "Error: POSTGRES_PASSWORD environment variable is not set"
    exit 1
fi

if [ -z "$POSTGRES_DB" ]; then
    echo "Error: POSTGRES_DB environment variable is not set"
    exit 1
fi

export PGPASSWORD=$POSTGRES_PASSWORD

if [ -z "$DO_NOT_START_CONTAINER" ]; then
    pid=$(finch ps -q -f "name=$CONTAINER_NAME")

    if [ "$pid" != "" ]; then
        echo "Stopping local DB... $pid"
        finch kill "$pid"
    fi

    pid=$(finch ps -q -a -f "name=$CONTAINER_NAME")

    if [ "$pid" != "" ]; then
        echo "Removing local DB... $pid"
        finch rm "$pid"
    fi

    echo "Starting local DB..."
    finch run \
        -e POSTGRES_USER=${POSTGRES_USER} \
        -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
        -e POSTGRES_DB=${POSTGRES_DB} \
        -p 6432:5432 \
        -d \
        --name $CONTAINER_NAME \
        postgres:14
fi

# Wait for the DB to start
while true; do
    psql -h localhost -p 6432 -U ${POSTGRES_USER} ${POSTGRES_DB} -c "SELECT 1;" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "Connected"
        break
    else
        echo "Connection failed. Retrying in 1 second..."
        sleep 1
    fi
done

psqldef -U ${POSTGRES_USER} -p 6432 -h localhost -f database/schema.sql ${POSTGRES_DB}
