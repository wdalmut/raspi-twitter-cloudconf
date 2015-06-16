default: all

all: test
		GOARM=6 GOARCH=arm GOOS=linux go build -a
		ssh pi@$(TARGET) 'killall picme | true'
		scp picme pi@$(TARGET):~
		scp config.json pi@$(TARGET):~
		ssh pi@$(TARGET) './picme < /dev/null >/tmp/picme.log 2>&1 &'

test:
		go test -v ./...

binary: all

