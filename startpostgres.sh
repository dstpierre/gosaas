#!/bin/bash

docker run --name postgres1 -p 5432:5432  -v ~/.pgdata:/var/lib/postgresql/data  -e POSTGRES_PASSWORD=dbpwd -d postgres
