all: build send exec
build:
	GOOS=linux GOARCH=arm GOARM=7 GO111MODULE=on go build
send:
	scp motionbot pi@pi:~/motionbot/
exec:
	ssh -t pi@pi "cd ~/motionbot; source config; ./motionbot"