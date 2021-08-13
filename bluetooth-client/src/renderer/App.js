import React from 'react';

import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';

import BluetoothIcon from '@material-ui/icons/Bluetooth';
import BluetoothDisabledIcon from '@material-ui/icons/BluetoothDisabled';

import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import Cell from './components/Cell';
import Chart from './components/Chart';
import UUID from './uuid';

export class App extends React.Component {

  constructor(props) {
    super(props);
  
    this.state = {
      connected: false,
    };

    this.device = null;
    this.connectDevice = this.connectDevice.bind(this);
    this.disconnectDevice = this.disconnectDevice.bind(this);
  };

  connectDevice() {
    navigator.bluetooth.requestDevice({
      filters: [{
        services: [UUID.SERVICE_ID]
      }]
    })
      .then(device => {        
        this.device = device;
        return device.gatt.connect();
      })
      .then(server => {
        this.serverId = server.device.id
        return server.getPrimaryService(UUID.SERVICE_ID);
      })
      .then(service => {
        return service.getCharacteristics();
      })
      .then(characteristics => {
        this.characteristicsList = characteristics;
        this.setState({connected: true});
        return;
      })
      .catch(error => {
        console.error(error);
      });    
    return;
  };

  disconnectDevice() {
    this.device.gatt.disconnect();
    this.setState({connected: false});
  }

  getCharacteristic(uuid) {
    return this.characteristicsList.find(i => i.uuid === uuid)
  };

  render() {
    let button;
    if (this.state.connected) {
      button = <Button
                  variant="contained"
                  color="secondary"
                  onClick={this.disconnectDevice}
                  startIcon={<BluetoothDisabledIcon />}>
                Disconnect
                </Button>
    }
    else {
      button = <Button
                   variant="contained"
                   color="primary"
                   onClick={this.connectDevice}
                   startIcon={<BluetoothIcon />}>
                Connect                
                </Button>
    }

    return (
      <Container maxWidth="sm">
        <Box my={4}>
          <Typography variant="h5" component="h1" gutterBottom>
            Bonzai Tree Broadcast
          </Typography>
          {button}
          {
            this.state.connected &&
              <div>
                <TableContainer component={Paper} style={{ width: 600 }}>
                  <Table aria-label="simple table">
                    <TableHead>
                      <TableRow>
                        <TableCell align="right">Soil Temperature (C°)</TableCell>
                        <TableCell align="right">Soil Moisture (C/cm^2)</TableCell>
                        <TableCell align="right">UV Index</TableCell>
                        <TableCell align="right">Humidity (%)</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      <TableRow>
                        <Cell
                          align="right"
                          characteristic={this.getCharacteristic(UUID.SOIL_TEMP)}>
                        </Cell>
                        <Cell
                          align="right"
                          characteristic={this.getCharacteristic(UUID.SOIL_MOISTURE)}>
                        </Cell>
                        <Cell
                          align="right"
                          characteristic={this.getCharacteristic(UUID.UV_INDEX)}>
                        </Cell>
                        <Cell
                          align="right"
                          characteristic={this.getCharacteristic(UUID.HUMIDITY)}>
                        </Cell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </TableContainer>
                <Chart characteristic={this.getCharacteristic(UUID.SOIL_MOISTURE)}
                          yLabel='Moisture (C/cm^2)'
                ></Chart>
                <Chart characteristic={this.getCharacteristic(UUID.UV_INDEX)}
                          yLabel='UV Index'
                ></Chart>
                <Chart characteristic={this.getCharacteristic(UUID.SOIL_TEMP)}
                          yLabel='Temperature (C°)'
                ></Chart>
              </div>
          }
        </Box>
      </Container>
    );
  }
}
