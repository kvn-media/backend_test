network:
	docker network create bank-network

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

serverdocker:
	docker run --name simple_bank -p 8080:8080 -e GIN_MODE=release --network bank-network backend_test:latest

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_postgres_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres backend_test

dropdb:
	docker exec -it postgres12 dropdb backend_test

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kvn-media/backend_test/db/sqlc Store

swagger:
	swag init -g ./api/server.go

dockercomposerebuild:
	docker compose up --force-recreate --build api

.PHONY: network postgres serverdocker createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock swagger dockercomposerebuild