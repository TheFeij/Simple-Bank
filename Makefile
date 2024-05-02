postgres:
	docker run --name postgres-container --network simple-bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:16rc1-alpine3.18

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

migrateup1:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

testmigrateup:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose up

testmigrateup1:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose up 1

testmigratedown:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose down

testmigratedown1:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/simple_bank_test?sslmode=disable" -verbose down 1

mockdb:
	mockgen -package mockdb -destination db/mock/services.go Simple-Bank/db/services Services

mocktokenmaker:
	mockgen -package mocktokenmaker -destination token/mock/maker.go Simple-Bank/token Maker

server:
	go run main.go

test:
	go test -v -cover ./...

proto:
	rm -f pb/*.go
	protoc --go_out=pb --proto_path=proto --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        proto/*.proto

.PHONY: postgres, createdb, dropdb, createtestdb, droptestdb, mockdb, mocktokenmaker
.PHONY: migratedown, migrateup, testmigratedown. testmigrateup, server
.PHONY: migratedown1, migrateup1, testmigrateup1, testmigratedown1, proto