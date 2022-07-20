package main

import (
	"fmt"
	"goGrpcConn/api/storage/postgres"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	bcr "goGrpcConn/api/core/blog"
	bgk "goGrpcConn/api/gunk/v1/admin/blog"
	bsr "goGrpcConn/api/service/blog"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error while configuration: %v", err)
	}
	store, err := dbConnection(config)
	if err != nil {
		log.Print("unable to configure storage", err)
	}
	serverPort := config.GetString("server.port")
	fmt.Printf("\n**********************************************************\n")
	fmt.Println("API service running on PORT : ", serverPort)
	fmt.Printf("**********************************************************\n\n")
	if err := setupGRPCService(store, config); err != nil {
		log.Printf("error lading configurtion: %v", err)

	}
}

func dbConnection(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupGRPCService(store *postgres.Storage, config *viper.Viper) error {
	if err := store.RunMigration(config.GetString("database.migrationDir")); err != nil {
		log.Printf("unable to run migrations")
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GetString("server.port")))
	if err != nil {
		log.Printf("Failed to listen on port %s: %v", config.GetString("server.port"), err)

	}

	grpcServer := grpc.NewServer()
	registerGrpcServices(grpcServer, store)
	log.Printf("Server api management listening at : %+v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Print("Failed to serve over Port: 8082")
		return err
	}

	return nil
}

func registerGrpcServices(grpcServer *grpc.Server, store *postgres.Storage) {
	bgk.RegisterBlogServiceServer(grpcServer, bsr.BlogCoreConn(bcr.ConnWithStorage(store)))
}
