.PHONY: build start stop check clean

tz-server: server/main.go
	go build -o ./tz-server server/main.go

build: tz-server

start: build
	./tz-server -port 8010 -tz America/Adak &
	./tz-server -port 8020 -tz America/Los_Angeles &
	./tz-server -port 8030 -tz Europe/Lisbon &
	./tz-server -port 8040 -tz Europe/Warsaw &
	./tz-server -port 8050 -tz Europe/Kyiv &
	./tz-server -port 8060 -tz Asia/Kolkata &

stop:
	killall tz-server

restart: stop build start

check:
	ps aux | grep tz-server

clean:
	rm ./tz-server

client: client/main.go
	go run client/main.go Alaska=8010 "Los Angeles=8020" Lisbon=8030 Warsaw=8040 Kyiv=8050 Kolkata=8060
