
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
	rm -rf app.o

rebuild: clean build