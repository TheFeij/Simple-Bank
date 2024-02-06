postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:16rc1-alpine3.18

createdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-container dropdb simple_bank

createtestdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank_test

droptestdb:
	docker exec -it postgres-container dropdb simple_bank_test

migrateup:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

testmigrateup:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose up

testmigratedown:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose down

mock:
	mockgen -package mockdb -destination db/mock/services.go Simple-Bank/db/services Services

server:
	go run main.go

test:
	go test -v -cover ./...

.PHONY: postgres, createdb, dropdb, createtestdb, droptestdb, mock
.PHONY: migratedown, migrateup, testmigratedown. testmigrateup, server