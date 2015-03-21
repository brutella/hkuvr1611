GO ?= go

run:
	$(GO) run daemon/hkuvr1611d.go

test:
	$(GO) test ./...

bbb:
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build daemon/hkuvr1611d.go

# 2014-03-21 Use GOARM=5 instead of GOARM=6 to prevent 'SIGILL: illegal instruction' error because of floating point issue see http://www.raspberrypi.org/forums/viewtopic.php?f=34&t=10781 - why does this happen after using it for so long?=
rpi:
	GOOS=linux GOARCH=arm GOARM=5 $(GO) build daemon/hkuvr1611d.go
