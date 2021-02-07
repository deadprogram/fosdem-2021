package main

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	must("enable BLE interface", adapter.Enable())

	// Start scanning.
	println("scanning...")
	must("start scan", adapter.Scan(scanHandler))
}

func scanHandler(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	println("device:", device.Address.String(), device.RSSI, device.LocalName())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
