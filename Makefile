.PHONY: build run clean profile-mem profile-cpu docker-build docker-run docker-stop docker-clean

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o go-minitrackr.exe ./cmd/server

run:
	go run ./cmd/server

clean:
	del /Q go-minitrackr.exe 2>nul || exit 0

docker-build:
	docker build -t go-minitrackr:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down -v
	docker rmi go-minitrackr:latest 2>nul || exit 0

profile-mem:
	go run ./cmd/server &
	timeout /t 5
	curl http://localhost:3000/debug/pprof/heap > heap.prof
	go tool pprof -http=:8080 heap.prof

profile-cpu:
	go run ./cmd/server &
	timeout /t 5
	curl http://localhost:3000/debug/pprof/profile?seconds=30 > cpu.prof
	go tool pprof -http=:8080 cpu.prof

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...
