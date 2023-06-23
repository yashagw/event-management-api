DB_URL=postgresql://root:secret@localhost:5443/event_management_db?sslmode=disable

postgres: 
	docker run --name postgres -p 5443:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root event_management_db

dropdb:
	docker exec -it postgres dropdb event_management_db

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

mock:
	mockgen -package mockdb -destination db/mock/mockdb.go github.com/yashagw/event-management-api/db Provider
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/yashagw/event-management-api/worker TaskDistributor

migratefile:
	migrate create -ext sql -dir db/migration -seq db_seq

create_admin:
	go build -o ./bin/create_admin ./scripts/create_admin.go
	chmod +x ./bin/create_admin
	./bin/create_admin

server:
	go run main.go

sql:
	make dropdb
	make createdb
	make migrateup