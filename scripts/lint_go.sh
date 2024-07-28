#!/bin/bash

for d in */ ; do
    if [ -f "$d/go.mod" ]; then
        echo "Linting $d"
        (
            cd $d || exit
            golangci-lint run
        ) 
    fi
done

wait 
