postgres:
	docker run --name db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234567890 -d postgres
createdb:
	docker exec -it db createdb -U root simple_bank
dropdb:
	docker exec -it db dropdb -U root simple_bank

create-migrate:
	migrate create -ext sql -dir db/migration/ -seq init 

migrateup:
	migrate -path db/migration -database "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run cmd/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/b0nbon1/bank-system/db/sqlc Store

.PHONY: createdb postgres dropdb migrate sqlc test server mock
