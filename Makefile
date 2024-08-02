up:
	docker compose up -d

exec:
	docker compose exec go bash

generate:
	@buf generate --path ./proto

start:
	go run ./cmd/*.go start

debug:
	CGO_ENABLED=0 go build -gcflags "all=-N -l" -buildvcs=false -o hello-debug ./cmd
	dlv --listen=:4000 --headless=true --api-version=2 exec ./hello-debug start

.PHONY: cert
cert:
	cd cert && ./gen.sh

