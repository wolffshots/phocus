![Worklow Status](https://github.com/wolffshots/phocus/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/wolffshots/phocus.svg)](https://pkg.go.dev/github.com/wolffshots/phocus/v2)
[![codecov](https://codecov.io/github/wolffshots/phocus/branch/main/graph/badge.svg?token=641UGV72AY)](https://codecov.io/github/wolffshots/phocus)

## phocus

A generic set of packages and app to speak to a device via serial and relay the responses to

Primarily built to communicate with Phocos branded inverters but the only thing you should need to change to use other inverters is the structure of the messages and the populating of the queue

### Installation

#### Ubuntu/systemd

If you don't need to differ from the default setup then it should be as simple as:

1.  Clone the repo

    ```sh
    git clone https://github.com/wolffshots/phocus.git
    ```

2.  Create a `config.json` file from `config.json.example` and customise your settings

3.  Run the install script (which builds the app then will ask for your password to put it in the correct place and makes a service file for it linked to `phocus.service` )

    ```sh
    cd phocus && ./install.sh
    ```

4.  (Re)start the service

    ```sh
    sudo service phocus restart
    ```

### Updating

To update you should just be able to pull/checkout the newer version, call `./install.sh` and restart the app with `sudo service phocus restart`
