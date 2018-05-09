#/bin/bash

docker run -it --rm --link postgres1:postgres postgres psql -h postgres -d gosaas -U postgres