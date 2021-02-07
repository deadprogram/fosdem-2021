package main

import (
	"fmt"
	"image/color"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/bluetooth"
	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyterm"
	"tinygo.org/x/tinyterm/fonts/proggy"
)

// replace this with the MAC address of the Bluetooth peripheral you want to connect to.
const deviceAddress = "E4:08:32:A8:9B:54"

var (
	display  st7789.Device
	terminal = tinyterm.NewTerminal(&display)

	black = color.RGBA{0, 0, 0, 255}
	font  = &proggy.TinySZ8pt7b

	adapter = bluetooth.DefaultAdapter
	device  *bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)
	buf     = make([]byte, 255)
)

func main() {
	initDisplay()

	must("enable BLE interface", adapter.Enable())

	println("start scan...")
	fmt.Fprintf(terminal, "\nstart scan...")

	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
		fmt.Fprintf(terminal, "\nconnected to %s", result.Address.String())
	}

	defer device.Disconnect()

	println("discovering services/characteristics")
	fmt.Fprintf(terminal, "\ndiscovering services/characteristics")

	srvcs, err := device.DiscoverServices(nil)
	must("discover services", err)

	for _, srvc := range srvcs {
		println("- service", srvc.UUID().String())
		fmt.Fprintf(terminal, "\n- service %s", srvc.UUID().String())

		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			println(err)
			fmt.Fprintf(terminal, "\n%s", err.Error())
		}
		for _, char := range chars {
			println("-- characteristic", char.UUID().String())
			n, err := char.Read(buf)
			if err != nil {
				println("    ", err.Error())
				fmt.Fprintf(terminal, "\n    %s", err.Error())
			} else {
				println("    data bytes", strconv.Itoa(n))
				fmt.Fprintf(terminal, "\n    data bytes %d", n)
				println("    value =", string(buf[:n]))
				fmt.Fprintf(terminal, "\n    value = %s", string(buf[:n]))
			}
		}
	}

	println("Done.")
}

func scanHandler(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	println("device:", device.Address.String(), device.RSSI, device.LocalName())
	fmt.Fprintf(terminal, "\ndevice: %s %d %s", device.Address.String(), device.RSSI, device.LocalName())
	if device.Address.String() == deviceAddress {
		adapter.StopScan()
		ch <- device
	}
}

func must(action string, err error) {
	if err != nil {
		for {
			println("failed to " + action + ": " + err.Error())
			time.Sleep(time.Second)
		}
	}
}

func initDisplay() {
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8000000,
		SCK:       machine.TFT_SCK,
		SDO:       machine.TFT_SDO,
		SDI:       machine.TFT_SDO,
		Mode:      0,
	})

	display = st7789.New(machine.SPI1,
		machine.TFT_RESET,
		machine.TFT_DC,
		machine.TFT_CS,
		machine.TFT_LITE)
	display.Configure(st7789.Config{
		Rotation:   st7789.ROTATION_180,
		Height:     320,
		FrameRate:  st7789.FRAMERATE_111,
		VSyncLines: st7789.MAX_VSYNC_SCANLINES,
	})
	display.FillScreen(black)

	terminal.Configure(&tinyterm.Config{
		Font:       font,
		FontHeight: 10,
		FontOffset: 6,
	})
}
