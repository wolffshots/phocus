App

![Worklow Status](https://github.com/wolffshots/phocus/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/wolffshots/phocus.svg)](https://pkg.go.dev/github.com/wolffshots/phocus)

# phocus

A generic set of packages and app to speak to a device via serial and relay the responses to

Primarily built to communicate with Phocos branded inverters but the only thing you should need to change to use other inverters is the structure of the messages and the populating of the queue

## Installation

Realistically you'll need to change some files to make it connect to your MQTT broker and serial device correctly but those changes should all be straightforward

### Ubuntu/systemd
If you don't need to differ from the default setup then it should be as simple as:

1. Clone the repo

    ```sh
    git clone https://github.com/wolffshots/phocus.git
    ```

2. Run the install script (which builds the app then will ask for your password to put it in the correct place and make a service file for it)

    ```sh
    cd phocus && ./install.sh
    ```

3. (Re)start the service

    ```sh
    sudo service phocus restart
    ```

## TODO
- [ ] explain config
- [ ] explain extra/drop in packages and using gomod replace - https://go.dev/doc/modules/gomod-ref
- [ ] separate queue population into phocus_messages
