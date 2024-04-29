migrateup:
	migrate -path db/postgres/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/postgres/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/postgres/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/postgres/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/postgres/migration -seq $(name)

sqlc:
	sqlc generate

.PHONY:
	migrateup migrateup1 migratedown migratedown1 new_migration sqlc