postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=david1380 -d postgres:16rc1-alpine3.18

createdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-container dropdb simple_bank

createtestdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank_test

droptestdb:
	docker exec -it postgres-container dropdb simple_bank_test

migrateup:
	 migrate -path ./db/migration -database "postgres://root:david1380@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path ./db/migration -database "postgres://root:david1380@localhost:5432/simple_bank?sslmode=disable" -verbose down


.PHONY: postgres, createdb, dropdb, createtestdb, droptestdb, migratedown, migrateup