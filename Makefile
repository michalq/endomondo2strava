.PHONY: all install run clean

install:
	mkdir -p ./tmp/workouts

run:
	go run cmd/sync/main.go

all: install run

clean:
	rm -rf ./tmp