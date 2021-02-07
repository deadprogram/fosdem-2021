module github.com/hybridgroup/fosdem2021/tinyflight

go 1.15

replace tinygo.org/x/tinyfont/freemono => ./fonts

require (
	gobot.io/x/gobot v1.15.0
	tinygo.org/x/bluetooth v0.2.1-0.20210201231738-27cc35a60b53
	tinygo.org/x/drivers v0.14.1-0.20210131100942-d6408374ed06
	tinygo.org/x/tinydraw v0.0.0-20200416172542-c30d6d84353c
	tinygo.org/x/tinyfont v0.2.1
	tinygo.org/x/tinyterm v0.0.0-20210131125732-23846e704bed
)
