#!/bin/bash

for d in */ ; do
    if [ -f "$d/go.mod" ]; then
        echo "Updating $d"
        (
            cd "$d" || exit
            go mod tidy
            go mod vendor
        ) &
    fi
done

wait