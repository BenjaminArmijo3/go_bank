makemigration:
	# creates a new migration
	goose -dir db/migrations create ${name} sql

migrate:
	goose -dir db/migrations up


db_up:
	docker exec -it bank_postgres createdb --username=root --owner=root bank_db

db_down:
	docker exec -it bank_postgres dropdb --username=root bank_db

sqlc:
	sqlc generate

start:
	CompileDaemon -command="./bank"

test:
	go test -v -cover ./...
