GO ?= go
	
test: 
	$(GO) test ./...

bbb:
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build daemon/hkuvr1611d.go

rpi:
	GOOS=linux GOARCH=arm GOARM=6 $(GO) build daemon/hkuvr1611d.go