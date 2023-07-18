#!/bin/bash

echo create dir and config
sudo mkdir -p /usr/bin/phocus_app
echo 
if test -e config.json; then
    echo config.json found
else
    echo config.json not found! generating from config.json.example
    cp config.json.example config.json
    read -n 1 -s -r -p "press any key to continue once you've checked that config.json is correct"
    echo
fi
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