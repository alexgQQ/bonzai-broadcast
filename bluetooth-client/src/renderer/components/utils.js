export function ab2str(buf) {
    return String.fromCharCode.apply(null, new Uint8Array(buf));
};

export function readCharacteristic(characteristic, setter) {
    characteristic.readValue()
        .then(i => i.buffer)
        .then(ab2str)
        .then((value) => {
            setter(value);
        });
};

export function subscribeCharacteristic(characteristic, setter) {
    characteristic.addEventListener('characteristicvaluechanged', event => {
        let value = ab2str(event.target.value.buffer);
        setter(value);
    });
    characteristic.startNotifications();
};

export function wait(sec) {
    return new Promise((resolve => {
        setTimeout(() => {
        resolve(true)
        }, 1000 * sec)
    }));
};