import React from 'react';
import { LineChart, Line, CartesianGrid, XAxis, YAxis } from 'recharts';
import { subscribeCharacteristic } from './utils';

function initData() {
    let data = [];
    for (let step = 0; step < 100; step++) {
        data.push({argument: step, value: 0});
    }
    return data;
};

function arrayRotate(arr, reverse) {
    if (reverse) arr.unshift(arr.pop());
    else arr.push(arr.shift());
    return arr;
}

export default class Chart extends React.Component {

    constructor(props) {
        super(props);
        this.dataSet = initData();
        this.state = {
            value: 0
        };
    }

    componentDidMount() {
        if (this.props.characteristic) {
            subscribeCharacteristic(this.props.characteristic, val => this.setState({value: val}));
        }
    }

    componentWillUnmount() {
        this.props.characteristic.removeEventListener('characteristicvaluechanged', () => {});
    }

    render() {

        this.dataSet[0].value = this.state.value;
        this.dataSet = arrayRotate(this.dataSet);
        this.dataSet.forEach(function(entry, index, arr) {
            arr[index].argument = index + 1
        });

        return (
            <LineChart width={600} height={200} data={this.dataSet.slice()}>
                <Line type="monotone" dataKey="value" stroke="#d92518" />
                <CartesianGrid stroke="#ccc" />
                <XAxis tick={false} axisLine={false} dataKey="argument" />
                <YAxis
                    tickLine={false} domain={['auto', 'auto']}
                    label={{ value: this.props.yLabel, angle: -90, position: 'insideBottomLeft' }} />
            </LineChart>
        );
    }
}
