DB_URL ?= postgresql://root:secret@192.168.29.20:5432/simple_bank?sslmode=disable


postgres17:
	fuser -k 5432/tcp 2>/dev/null || true && docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

wait-for-db:
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec postgres17 pg_isready -h localhost -U root >/dev/null 2>&1; do \
		sleep 1; \
	done
	@sleep 2
	@echo "PostgreSQL is fully ready!"

psqldrop:
	docker stop postgres17
	docker rm postgres17

createdb:
	docker exec -it postgres17 createdb -h localhost --username=root --owner=root simple_bank

# New target: Kill all connections to simple_bank
killconnections:
	docker exec -i postgres17 psql -U root -d postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'simple_bank' AND pid <> pg_backend_pid();"

# Updated dropdb: kill connections first, then drop
dropdb: killconnections
	docker exec -it postgres17 dropdb simple_bank

# Alternative: Drop with force (PostgreSQL 13+)
dropdbforce:
	docker exec -i postgres17 psql -U root -d postgres -c "DROP DATABASE simple_bank WITH (FORCE);"

# MIGRATION
# ********************************************************************************
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migratefresh: migratedown migrateup
	@echo "Fresh migration complete"

truncate:
	docker exec -i postgres17 psql -U root -d simple_bank -c "TRUNCATE entries, transfers, accounts RESTART IDENTITY CASCADE;"
# ********************************************************************************

sqlc:
	sqlc generate

# MIGRATION
# ********************************************************************************
testconnection:
	go test -v ./db/tests/main_test.go

testoverall:
	go test -v -cover -coverpkg=github.com/suryansh74/simplebank/db/sqlc -count=1 ./db/tests

testapi:
	go test -v -coverprofile=coverage.out ./api

testutil:
	go test -v -cover ./api
# ********************************************************************************

server:
	fuser -k 8000/tcp 2>/dev/null || true && go run .

mock:
	mockgen -source=db/store.go -destination=db/mock/store.go -package=mock Store

scratch: postgres17 wait-for-db createdb migrateup testconnection testoverall testapi testutil

.PHONY: postgres17 createdb dropdb dropdbforce killconnections migrateup migratedown sqlc testconnection testoverall psqldrop server mock
