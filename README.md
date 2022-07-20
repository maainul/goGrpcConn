
# START PROJECT :
    docker start postgres12
    go run .
# 1. Install Postgres Docker images

1. Pull image

    docker pull postgres:12-alpine

2. Start a postgres instance

    docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

3. Try to connect and exect its console

    docker exec -it postgres12 psql -U root

    root=# select now();
    
    \q

4. View Container logs

    dokcer logs postgres12

5. Go to postgresdb

    docker exec -it postgres12 /bin/sh

        pwd

        ls -l

        createdb --username=root --owner=root simpe_bank 

        psql simple_bank

        dropdb simple_bank

        exit

6. Crate db from outside of the container

        docker exec -it postgres12 createdb --username=root --owner=root simple_bank

7. Excess Database from command

        docker exec -it postgres12 psql -U root simple_bank


8. Create Makefile

9. Stop and remove container

    docker stop postgres12

    docker ps

    docker rm postgres12

    docker ps -a

10. Run makefile

    make postgres

    docker ps

    make createdb

    
## 2. INSTALL migrations/CLI for Migrations

    curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
   
    sudo apt-get update
    
    sudo apt-get install -y migrate

### Run the migration

### Create migration files

    migrate --help

    migrate create -ext sql -dir db/migrations -seq migration_file_name

Migration UP :

    migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/test_bank?sslmode=disable" -verbose up

Migration DOWN :

    migrate -path db/migrations/ -database "postgresql://root:secret@localhost:5432/test_bank?sslmode=disable" -verbose down




Details in this blog:

https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g

## UNIT TEST:
Details in this blog
https://dev.to/techschoolguru/write-go-unit-tests-for-db-crud-with-random-data-53no
