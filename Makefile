# generate migrations, $(name) - name of the migration
generate_migrations:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path internal/db/migrations -database "postgresql://devuser:admin@localhost:5432/cv_db?sslmode=disable" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path internal/db/migrations -database "postgresql://devuser:admin@localhost:5432/cv_db?sslmode=disable" -verbose down

# generate db related go code with sqlc
# for windows:	cmd.exe /c "docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate"
sqlc:
	sqlc generate

# run all tests
test:
	go test -v -cover ./...

# run tests in the given path (p) and display results in the html file
test_coverage:
	go test $(p) -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: generate_migrations, migrate_up, migrate_down, sqlc