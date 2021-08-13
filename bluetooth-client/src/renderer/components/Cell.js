import React, { useState, useEffect } from 'react';
import TableCell from '@material-ui/core/TableCell';
import { subscribeCharacteristic } from './utils';

export default function Cell(props) {

    const [value, setValue] = useState('');

    useEffect(() => {
      if (props.characteristic) {
          subscribeCharacteristic(props.characteristic, setValue);
          return () => {
              props.characteristic.removeEventListener('characteristicvaluechanged', () => {});
          }
      }
    });    
    return (
      <TableCell {...props}>{value}</TableCell>
    );
};