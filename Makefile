# -------------------------
#	Migration
# -------------------------

# Install postgres migration tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Migrate UP
migrate:
	migrate -source file://configs -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" up

# Migrate DROP
migrate-drop:
	migrate -source file://configs -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" drop

# -------------------------
#	App build & run
# -------------------------

# For next targets
app.o:
	go build -o app.o ./cmd/server/...

# Build executable
build: app.o

# Run executable
run: app.o
	./app.o

# Clean compiled & generated files
clean:
	rm -rf app.o tests/*.txt

# Rebuild executable
rebuild: clean build

# Rebuild and run executable
rerun: rebuild run

# Run load test via Apache Bench
load-test:
	ab -n 1000 -c 100 "localhost:8080/user_banner?tag_id=365172912&feature_id=2130129125" > tests/user_banner_rps.txt