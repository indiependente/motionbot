all: build send
build:
	GOOS=linux GOARCH=arm GOARM=7 GO111MODULE=on go build
send:
	scp motionbot pi@pi:~/motionbot/