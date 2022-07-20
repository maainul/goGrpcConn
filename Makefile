DB_URL=postgresql://root:secret@localhost:5434/pfdb?sslmode=disable

postgres:
	docker run --name postgres12 -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root pfdb

ds:
	docker start postgres12

dropdb:
	docker exec -it postgres12 pfdb

server:
	go run main.go

test:
	go test -v --cover ./storage/postgres

.PHONY: postgres createdb dropdb server test ds
