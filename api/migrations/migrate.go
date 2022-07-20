package main

import (
	"log"
	"goGrpcConn/api/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
