.PHONY: build start stop check clean

tz-server: server/main.go
	go build -o ./tz-server server/main.go

build: tz-server

start: build
	./tz-server -port 8010 -tz America/Adak &
	./tz-server -port 8020 -tz Europe/London &
	./tz-server -port 8030 -tz Europe/Warsaw &
	./tz-server -port 8040 -tz Europe/Kyiv &

stop:
	killall tz-server

restart: stop build start

check:
	ps aux | grep tz-server

clean:
	rm ./tz-server
