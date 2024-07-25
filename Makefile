up:
	docker compose up -d

exec:
	docker compose exec go bash

generate:
	@buf generate --path ./proto
