make:
	GOARM=7 GOARCH=arm GOOS=linux go build

fmt:
	gofmt -w advertiser.go

deps:
	go get gopkg.in/yaml.v2
	go get -d -u gobot.io/x/gobot/...

install:
	scp advertiser pi@${LOGGER_IP}:~

run:
	ssh pi@${LOGGER_IP} sudo ./advertiser

ssh:
	ssh pi@${LOGGER_IP}
