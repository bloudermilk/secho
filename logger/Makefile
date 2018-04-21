make:
	GOARM=7 GOARCH=arm GOOS=linux go build

deps:
	go get gopkg.in/yaml.v2
	go get -d -u gobot.io/x/gobot/...

install:
	scp logger config.example.yml pi@${LOGGER_IP}:~

run:
	ssh pi@${LOGGER_IP} ./logger config.example.yml

ssh:
	ssh pi@${LOGGER_IP}
