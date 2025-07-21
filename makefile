CLIENT_DIR := web-ui
VITE := npx vite

# start db 
postgres:
	docker run --name db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234567890 -d postgres

# create db 
createdb:
	docker exec -it postgresdb createdb -U root stratal

# drop the db
dropdb:
	docker exec -it db dropdb -U root stratal

# create new migration for db
create-migrate:
	migrate create -ext sql -dir internal/storage/db/migration/ -seq init
# create new migration for db with a name
create-migrate-with-name:
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "Migration name cannot be empty"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir internal/storage/db/migration/ -seq "$$name"

# run migrations
migrateup:
	migrate -path internal/storage/db/migration -database "postgresql://root:1234567890@localhost:5432/stratal?sslmode=disable" -verbose up

# spin down migrations
migratedown:
	migrate -path internal/storage/db/migration -database "postgresql://root:1234567890@localhost:5432/stratal?sslmode=disable" -verbose down

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

build_worker:
	env GOOS=linux CGO_ENABLED=0 go build -o bin/worker cmd/worker/main.go

run_worker_dev:
	go run cmd/worker/main.go
	
run_server_dev:
	go run cmd/server/main.go

# create mockdb for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/b0nbon1/bank-system/db/sqlc Store

.PHONY: createdb postgres dropdb migrate sqlc test server mock front_dev clean_front build_front
