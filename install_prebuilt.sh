#!/bin/bash

echo create dir and config
sudo mkdir -p /usr/bin/phocus_app
echo 
if test -e config.json; then
    echo config.json found, will be moved to app dir
else
    echo config.json not found! generating from latest config.json.example on github
    wget https://github.com/wolffshots/phocus/releases/latest/download/config.json.example
    mv config.json.example config.json
    read -n 1 -s -r -p "press any key to continue once you've checked that config.json (in the current directory) is correct"
    echo
fi
sudo mv config.json /usr/bin/phocus_app/
echo done
echo 
echo pulling latest binary and service from github
wget https://github.com/wolffshots/phocus/releases/latest/download/phocus.service
wget https://github.com/wolffshots/phocus/releases/latest/download/phocus
echo done
echo 
echo move executable and service to dir
sudo mv ./phocus.service /usr/bin/phocus_app
sudo mv ./phocus /usr/bin/phocus_app
echo done
echo 
echo "link service (in app dir) to systemd dir"
sudo ln -f ./usr/bin/phocus_app/phocus.service /usr/lib/systemd/system/phocus.service
echo done