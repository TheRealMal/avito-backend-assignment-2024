
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate:
	migrate -source file://configs -database "${DB_URL}" up

app.o: 
	go build -o app.o ./cmd/server/...

build: app.o

run: app.o
	./app.o

clean:
	rm -rf app.o tests/*.txt

rebuild: clean build

load-test:
	ab -n 1000 -c 100 -g report "http://localhost:8080/user_banner?tag_id=1&feature_id=1" > tests/user_banner_rps.txt