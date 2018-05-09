#!/bin/bash

clear
rm gosaas-dev
go build -o gosaas-dev
./gosaas-dev -driver postgres -datasource "postgres://postgres:dbpwd@localhost/gosaas?sslmode=disable"