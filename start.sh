#!/bin/bash

clear
rm gosaas-dev
go build -o gosaas-dev
./gosaas-dev -driver mongo -datasource "127.0.0.1"