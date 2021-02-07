package main

import (
	"math/rand"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter

	heartRateMeasurement bluetooth.Characteristic
	heartRate            uint8 = 75 // 75bpm
)

func main() {
	must("enable BLE adaptor", adapter.Enable())

	adv := adapter.DefaultAdvertisement()
	must("config advertisement", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Go HRS",
		ServiceUUIDs: []bluetooth.UUID{bluetooth.ServiceUUIDHeartRate},
	}))
	must("start advertising", adv.Start())

	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHeartRate,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &heartRateMeasurement,
				UUID:   bluetooth.CharacteristicUUIDHeartRateMeasurement,
				Value:  []byte{0, heartRate},
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))

	startHeartbeat()
}

func startHeartbeat() {
	nextBeat := time.Now()
	for {
		nextBeat = nextBeat.Add(time.Minute / time.Duration(heartRate))
		println("tick", time.Now().Format("04:05.000"))
		time.Sleep(nextBeat.Sub(time.Now()))

		// random variation in heartrate
		heartRate = randomInt(65, 85)

		// and push the next notification
		heartRateMeasurement.Write([]byte{0, heartRate})
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

// Returns an int >= min, < max
func randomInt(min, max int) uint8 {
	return uint8(min + rand.Intn(max-min))
}
