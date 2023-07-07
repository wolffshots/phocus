#!/bin/bash

echo create dir and config
sudo mkdir -p /usr/bin/phocus_app
sudo cp config.json /usr/bin/phocus_app/
echo done
echo 
echo build app
go build -o ./build/phocus main.go
echo done
echo 
echo move executable to dir
sudo mv ./build/phocus /usr/bin/phocus_app
echo done
echo 
echo link service to systemd dir
sudo ln -f ./phocus.service /usr/lib/systemd/system/phocus.service
echo done