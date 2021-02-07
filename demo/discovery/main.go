package main

import (
	"os"
	"strconv"

	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  *bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)
	buf     = make([]byte, 255)
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
	srvcs, err := device.DiscoverServices(nil)
	must("discover services", err)

	for _, srvc := range srvcs {
		println("- service", srvc.UUID().String())

		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			println(err)
		}
		for _, char := range chars {
			println("-- characteristic", char.UUID().String())
			n, err := char.Read(buf)
			if err != nil {
				println("    ", err.Error())
			} else {
				println("    data bytes", strconv.Itoa(n))
				println("    value =", string(buf[:n]))
			}
		}
	}

	println("Done.")
}

func scanHandler(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
	println("found device:", result.Address.String(), result.RSSI, result.LocalName())
	if result.Address.String() == connectAddress() {
		adapter.StopScan()
		ch <- result
	}
}

func connectAddress() string {
	if len(os.Args) < 2 {
		println("usage: discover [address]")
		os.Exit(1)
	}

	// look for device with specific address
	address := os.Args[1]

	return address
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
