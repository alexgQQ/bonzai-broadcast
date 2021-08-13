# Bonsai Plant Monitor

This is a collection of software for a little side project that acts as a plant health monitor.
The general idea is to log important plant health values via a set of sensors, namely:
    * Air Humidity
    * Ambient Temperature
    * UV Light Level
    * Soil Moisture

The software in this application consist of:
    * sensor-drivers - A python script to log sensor values to a csv file that can be installed as a unix daemon.
    * ble-gatt-server - A bluetooth server written in go to relay sensor values that can be installed as a unix daemon.
    * term-dashboard - A terminal ui dashboard written in go to view logged sensor data.
    * bluetooth-client - A client application using react and electron to display sensor data over bluetooth.

## Build

Set golang environment to compile the binaries for an RPI3 architecture.
If using this with a RPI-Zero then `GOARM=6`

```bash
GOOS=linux GOARCH=arm GOARM=7
go build -o ble-gatt-server/bin/bonsaiServer ble-gatt-server/src/server.go
go build -o term-dashboard/dashboard term-dashboard/dashboard.go
```

## Deploy

Sync files between the dev environment and the device.

```bash
rsync -aiz -e ssh sensor-drivers term-dashboard ble-gatt-server rpi:~/bonsai-broadcast --delete
```

## Install

```bash
sudo ./sensor-drivers/install.sh
sudo ./ble-gatt-server/install.sh
```
