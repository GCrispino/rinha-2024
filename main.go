package main

import (
	"fmt"
	"os"
	"strings"

	database "github.com/GCrispino/rinha/internal/database/connection"
	"github.com/GCrispino/rinha/internal/database/repository"
	"github.com/GCrispino/rinha/internal/server"
	"github.com/GCrispino/rinha/internal/usecases/customers"
)

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

	driverName := "postgres"
	dbConnStr := "postgres://user:password@localhost/rinha?sslmode=disable"
	dbConn, err := database.NewDBConn(driverName, dbConnStr)
	if err != nil {
		panic(err)
	}

    customersRepo := repository.NewCustomers(dbConn)
    customersUsecase := customers.NewCustomerUsecase(customersRepo)

	s := server.NewServer(customersUsecase)

	s.Logger.Fatal(s.Start(bindAddr))
}
