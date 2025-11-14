DB_URL ?= postgresql://root:secret@192.168.29.20:5432/simple_bank?sslmode=disable

postgres17:
	fuser -k 5432/tcp 2>/dev/null || true && docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

psqldrop:
	docker stop postgres17
	docker rm postgres17

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

truncate:
	docker exec -i postgres17 psql -U root -d simple_bank -c "TRUNCATE entries, transfers, accounts RESTART IDENTITY CASCADE;"

sqlc:
	sqlc generate

testconnection:
	go test -v ./db/tests/main_test.go

testoverall:
	go test -v -cover -coverpkg=github.com/suryansh74/simplebank/db/sqlc -count=1 ./db/tests

server:
	fuser -k 8000/tcp 2>/dev/null || true && go run .                                                                                                                                        ─╯

.PHONY: postgres17 createdb dropdb migrateup migratedown sqlc testconnection testoverall psqldrop server
