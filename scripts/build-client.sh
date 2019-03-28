#!/usr/bin/env bash

docker run -it --rm \
    -v $(pwd)/client:/usr/src/app \
    -u $(id -u):$(id -g) \
    node:8 \
    bash -c "cd /usr/src/app && npm install && npm run build"
