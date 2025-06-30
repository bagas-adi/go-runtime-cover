#!/bin/bash

go build -o bin/pkg_injector cmd/pkg_injector/main.go; chmod +x bin/pkg_injector
bin/pkg_injector -filename=examples/integration/main.go