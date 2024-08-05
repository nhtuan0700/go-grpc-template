up:
	docker compose up -d

exec:
	docker compose exec go bash

generate:
	@buf generate --path ./proto

start:
	go run ./cmd/*.go start

debug:
	/go/bin/dlv --listen=:4000 --headless=true --log=true --accept-multiclient --api-version=2 exec ./cmd/*.go start
