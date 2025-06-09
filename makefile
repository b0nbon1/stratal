CLIENT_DIR := client
VITE := npx vite

# start db 
postgres:
	docker run --name db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234567890 -d postgres

# create db 
createdb:
	docker exec -it db createdb -U root stratal

# drop the db
dropdb:
	docker exec -it db dropdb -U root stratal

# create new migration for db
create-migrate:
	migrate create -ext sql -dir db/migration/ -seq init 

# run migrations
migrateup:
	migrate -path db/migration -database "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable" -verbose up

# spin down migrations
migratedown:
	migrate -path db/migration -database "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable" -verbose down

# generate new sqlc data models changes
sqlc:
	sqlc generate

# Run tests
test:
	go test -v -cover ./...

# run the server
server:
	go run cmd/main.go

# run development front-end
front_dev:
	cd $(CLIENT_DIR) && $(VITE)

# build front-end
build_front:
	cd $(CLIENT_DIR) && npm run build

clean_front:
	cd $(CLIENT_DIR) && rm -rf dist


# create mockdb for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/b0nbon1/bank-system/db/sqlc Store

.PHONY: createdb postgres dropdb migrate sqlc test server mock front_dev clean_front build_front
