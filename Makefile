

default: all

all: test
		GOARM=6 GOARCH=arm GOOS=linux go build -a
		ssh pi@192.168.1.16 'killall picme | true'
		scp start.sh pi@192.168.1.16:~
		scp picme pi@192.168.1.16:~
		scp config.json pi@192.168.1.16:~
		ssh pi@192.168.1.16 './picme < /dev/null >/tmp/picme.log 2>&1 &'

test:
		go test -v ./...

binary: all

