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


# Docker compose build & run
run:
	docker-compose up -d --build

# Docker compose stop
stop:
	docker-compose down

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
	golangci-lint run -c .golangci.yml

# Clean compiled & generated files
clean:
	rm -rf test/*

# Run load test via Apache Bench
load-test:
	go run ./cmd/token/... ADMIN > test/.token
	curl -s \
	-H "Authorization: Token $$(<test/.token)" \
	"localhost:8080/banner" | \
	python3 -c "import sys, json; print(json.load(sys.stdin)[-1]['feature'])" > test/.feat
	curl -s \
	-H "Authorization: Token $$(<test/.token)" \
	"localhost:8080/banner" | \
	python3 -c "import sys, json; print(json.load(sys.stdin)[-1]['tags'][-1])" > test/.tag
	
	ab -n 1000 -c 100 \
	-H "Authorization: Token $$(<test/.token)" \
	"localhost:8080/user_banner?tag_id=$(<.tag)&feature_id=$(<.feat)&use_last_revision=true" \
	> test/load_user_banner_straight_to_db.txt

	ab -n 1000 -c 100 \
	-H "Authorization: Token $$(<test/.token)" \
	"localhost:8080/user_banner?tag_id=$(<.tag)&feature_id=$(<.feat)&use_last_revision=false" \
	> test/load_user_banner_cached.txt

	ab -n 1000 -c 100 \
	-H "Authorization: Token $$(<test/.token)" \
	"localhost:8080/banner" \
	> test/load_banner.txt

	rm test/.feat test/.tag test/.token