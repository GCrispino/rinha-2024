package main

import (
	"fmt"
	"os"
	"strings"

	database "github.com/GCrispino/rinha-2024/internal/database/connection"
	"github.com/GCrispino/rinha-2024/internal/database/repository"
	"github.com/GCrispino/rinha-2024/internal/server"
	"github.com/GCrispino/rinha-2024/internal/usecases/customers"
)

type config struct {
	dbConnStr string
}

func loadConfig() config {
	return config{
		dbConnStr: os.Getenv("DB_CONN_STR"),
	}
}

func main() {
	args := os.Args
	lArgs := len(args)
	if lArgs > 2 {
		fmt.Println("USAGE: ./rinha <port>")
		os.Exit(1)
	}

	bindAddr := "8080"
	if lArgs == 2 {
		bindAddr = args[1]
	}
	if !strings.HasPrefix(bindAddr, ":") {
		bindAddr = ":" + bindAddr
	}

	cfg := loadConfig()

	driverName := "postgres"
	dbConn, err := database.NewDBConn(driverName, cfg.dbConnStr)
	if err != nil {
		panic(err)
	}
	defer dbConn.Conn.Close()

	customersRepo := repository.NewCustomers(dbConn)
	customersUsecase := customers.NewCustomerUsecase(customersRepo)

	s := server.NewServer(customersUsecase)

	s.Logger.Fatal(s.Start(bindAddr))
}
