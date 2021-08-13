package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"os/exec"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
)

type beacon struct {
	uuid ble.UUID
	command string
}

var beacons = []beacon{
	beacon{ // Device Model
		uuid: ble.MustParse("00010000-0002-1000-8000-00805F9B34FB"),
		command: "cat /proc/device-tree/model",
	},
	beacon{ // CPU Temp
		uuid: ble.MustParse("00010000-0003-1000-8000-00805F9B34FB"),
		command: "vcgencmd measure_temp | cut -d = -f 2",
	},
	beacon{ // CPU Load
		uuid: ble.MustParse("00010000-0004-1000-8000-00805F9B34FB"),
		command: "top -bn1 | grep load | awk '{printf \"%.2f%%\", $(NF-2)}'",
	},
	beacon{ // Memory Usage
		uuid: ble.MustParse("00010000-0005-1000-8000-00805F9B34FB"),
		command: "free -m | awk 'NR==2{printf \"%s/%sMB\", $3,$2 }'",
	},
	beacon{ // Uptime
		uuid: ble.MustParse("00010000-0006-1000-8000-00805F9B34FB"),
		command: "uptime -p | cut -d 'p' -f 2 | awk '{ printf \"%s\", $0 }'",
	},
	beacon{ // Soil Temperature
		uuid: ble.MustParse("00010000-0007-1000-8000-00805F9B34FB"),
		command: "tail -n 1 /home/pi/data/data.csv | awk -F, '{print $2}'",
	},
	beacon{ // Soil Moisture
		uuid: ble.MustParse("00010000-0008-1000-8000-00805F9B34FB"),
		command: "tail -n 1 /home/pi/data/data.csv | awk -F, '{print $3}'",
	},
	beacon{ // UV Index
		uuid: ble.MustParse("00010000-0009-1000-8000-00805F9B34FB"),
		command: "tail -n 1 /home/pi/data/data.csv | awk -F, '{print $6}'",
	},
	beacon{ // Ambient Humidity
		uuid: ble.MustParse("00010000-000A-1000-8000-00805F9B34FB"),
		command: "tail -n 1 /home/pi/data/data.csv | awk -F, '{print $8}'",
	},
}

var (
	TestSvcUUID = ble.MustParse("00010000-0001-1000-8000-00805F9B34FB")
	duration = 10 * time.Second
)

func CharFactory(uuid ble.UUID, cmd_string string) *ble.Characteristic {
	c := ble.NewCharacteristic(uuid)
	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		cmd := exec.Command("bash", "-c", cmd_string)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf( "Error: %s", err )
		} else {
			log.Printf("Output: %s", output)
			fmt.Fprintf(rsp, "%s", output)
		}
		fmt.Printf("%s\n", output)
	}))
	c.HandleNotify(ble.NotifyHandlerFunc(func(req ble.Request, n ble.Notifier) {
		log.Printf("count: Notification subscribed")
		for {
			select {
			case <-n.Context().Done():
				log.Printf("count: Notification unsubscribed")
				return
			case <-time.After(time.Second):
				cmd := exec.Command("bash", "-c", cmd_string)
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Printf( "Error: %s", err )
				} else {
					log.Printf("Output: %s", output)
				}
				if _, err := fmt.Fprintf(n, "%s", output); err != nil {
					log.Printf("count: Failed to notify : %s", err)
					return
				}
			}
		}
	}))
	return c
}

func main() {

	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("can't create new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	testSvc := ble.NewService(TestSvcUUID)
	for i, s := range beacons {
		testSvc.AddCharacteristic(CharFactory(s.uuid, s.command))
		i++
	}

	if err := ble.AddService(testSvc); err != nil {
		log.Fatalf("can't add service: %s", err)
	}

	ctx := ble.WithSigHandler(context.WithCancel(context.Background()))
	chkErr(ble.AdvertiseNameAndServices(ctx, "Gopher", testSvc.UUID))
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
