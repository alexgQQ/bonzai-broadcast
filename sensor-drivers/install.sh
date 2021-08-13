#!/bin/bash

SERVICE_NAME=bonsaiSensor.service
EXECUTABLE=bonsaiSensor

__usage="
Usage: ./install.sh [OPTIONS]

Install the plant monitor sensor script and dependencies as a service.
This must be run with elevated priviledges

Options:
  -u, --uninstall              Remove executable and service from the system
  -h, --help                   Show this help screen
  -r, --reload                 Replace the executable and reload the service
  --no-deps                    Install the service without dependencies
"

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

install_deps () {
    apt-get update -yq
    apt-get upgrade -yq
    apt install -yq libgpiod2 python3 python3-pip

    pip3 -q install adafruit-blinka \
                    adafruit-circuitpython-seesaw \
                    adafruit-circuitpython-dht \
                    adafruit-circuitpython-mcp3xxx
}

install_service () {
    # Move executable and service config to common locations
    cp ./$SERVICE_NAME /lib/systemd/system
    cp ./$EXECUTABLE /usr/bin

    # Make directory for data file
    mkdir -p /home/pi/data

    # Register and begin the new service
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    systemctl start $SERVICE_NAME
}

uninstall_service () {
    # Stop service and detach
    systemctl stop $SERVICE_NAME
    systemctl disable $SERVICE_NAME

    # Remove service file and executable
    rm /lib/systemd/system/$SERVICE_NAME
    rm /usr/bin/$EXECUTABLE
    systemctl daemon-reload
}

case "$1" in
    -u|--uninstall) # Remove service file and executable
        uninstall_service
        exit 0
        ;;
    -r|--reload) # Replace service executable
        # Stop service and detach
        systemctl stop $SERVICE_NAME

        # Remove service file and executable
        rm /usr/bin/$EXECUTABLE
        cp ./$EXECUTABLE /usr/bin

        systemctl start $SERVICE_NAME
        exit 0
        ;;
    --no-deps) # Install the service alone
        install_service
        exit 0
        ;;
    -h|--help) # Show the usage help text
        echo "$__usage"
        exit 0
        ;;
    *) # Install dependencies and service by default
        install_deps
        install_service
        exit 0
        ;;
esac
