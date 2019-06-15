#!/bin/bash

docker run -t --rm -v "$PWD":/usr/src/app -p "4000:4000" starefossen/github-pages