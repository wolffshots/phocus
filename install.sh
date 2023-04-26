#!/bin/bash

go build -o ./build/phocus main.go
sudo mv ./build/phocus /usr/bin
sudo ln -f ./phocus.service /usr/lib/systemd/system/phocus.service
