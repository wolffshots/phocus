#!/usr/bin/bash

go build -o /usr/bin/phocus main.go
ln -f ./phocus.service /usr/lib/systemd/system/phocus.service
