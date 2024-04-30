run:
	go run cmd/http/main.go

build:
	@go build -o bin/go-boilerplate cmd/http/main.go

test:
	@go test -v ./...

sqlc:
	sqlc generate

migrateup:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/postgres/migration -seq $(name)

.PHONY:
	run build test sqlc migrateup migrateup1 migratedown migratedown1 new_migration 