scan:
	cd ./demo/scanner; go run .; cd ../..

scanner-ble:
	cd ./demo/scanner-ble; tinygo flash -size short -target=clue .; cd ../..

discovery:
	cd ./demo/discovery; go run . ${ADDRESS}; cd ../..

discovery-clue:
	cd ./demo/discovery-clue; tinygo flash -size short -target=clue .; cd ../..

advertising:
	cd ./demo/advertising; tinygo flash -size short -target=circuitplay-bluefruit .; cd ../..

heartrate:
	cd ./demo/heartrate; tinygo flash -size short -target=itsybitsy-nrf52840 .; cd ../..

heartrate-monitor:
	cd ./demo/heartrate-monitor/; go run . ${ADDRESS}; cd ../..

tinyflight:
	cd ./demo/tinyflight; tinygo flash -size short -target=itsybitsy-nrf52840 .; cd ../..
