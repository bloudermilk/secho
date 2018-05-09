make:
	GOARM=7 GOARCH=arm GOOS=linux go build

deps:
	go get ./...

install:
	scp secho ../../examples/config.example.yml pi@${REMOTE_IP}:~

run:
	ssh pi@${REMOTE_IP} sudo ./secho config.example.yml

ssh:
	ssh pi@${REMOTE_IP}
