package main

import (
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  *bluetooth.Device

	heartRateServiceUUID        = bluetooth.ServiceUUIDHeartRate
	heartRateCharacteristicUUID = bluetooth.CharacteristicUUIDHeartRateMeasurement

	ch = make(chan bluetooth.ScanResult, 1)
)

func main() {
	must("enable BLE interface", adapter.Enable())

	println("scanning...")
	must("scan for specific peripheral", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
	}

	defer device.Disconnect()

	println("discovering services/characteristics")
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{heartRateServiceUUID})
	if err != nil || len(srvcs) == 0 {
		panic("could not find heart rate service")
	}

	srvc := srvcs[0]
	println("found service", srvc.UUID().String())

	chars, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{heartRateCharacteristicUUID})
	if err != nil || len(chars) == 0 {
		panic("could not find heart rate characteristic")
	}

	char := chars[0]
	println("found heart rate characteristic", char.UUID().String())

	char.EnableNotifications(func(buf []byte) {
		println("data:", uint8(buf[1]))
	})

	select {}
}

func scanHandler(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
	println("found device:", result.Address.String(), result.RSSI, result.LocalName())
	if result.Address.String() == connectAddress() {
		adapter.StopScan()
		ch <- result
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
