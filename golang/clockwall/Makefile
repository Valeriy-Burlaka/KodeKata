.PHONY: build start stop check clean

tz-server: server/main.go
	go build -o ./tz-server server/main.go

build: tz-server

start: build
	TZ=America/Adak ./tz-server -port 8010 &
	TZ=Europe/London ./tz-server -port 8020 &
	TZ=Europe/Warsaw ./tz-server -port 8030 &
	TZ=Europe/Kyiv ./tz-server -port 8040 &

stop:
	killall tz-server

restart: stop build start

check:
	ps aux | grep tz-server

clean:
	rm ./tz-server
