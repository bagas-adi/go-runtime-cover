#!/bin/bash

timestamp=$(date "+%Y-%m-%d %H%M%S")
coverage_dir=$(ls -td -- ./coverage-data/*/ | head -n 1)
echo "Coverage data directory: $coverage_dir"

finish_report(){
    go tool covdata textfmt -i=$coverage_dir -o $coverage_dir/coverage.out
    go tool cover -func=$coverage_dir/coverage.out
    exit 1
}

trap finish_report SIGINT
while true; do
    echo "[$timestamp]: ($coverage_dir) Waiting for coverage data..."
    if [ -d "$coverage_dir" ]; then
        go tool covdata percent -i=$coverage_dir  
        timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    fi
    sleep 1
done
