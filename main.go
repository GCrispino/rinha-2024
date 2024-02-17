package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "go.uber.org/automaxprocs"

	database "github.com/GCrispino/rinha-2024/internal/database/connection"
	"github.com/GCrispino/rinha-2024/internal/database/repository"
	"github.com/GCrispino/rinha-2024/internal/server"
	"github.com/GCrispino/rinha-2024/internal/usecases/customers"
)

type config struct {
	dbConnStr      string
	dbMaxOpenConns int
}

func loadConfig() config {
	dbConnStr := "postgres://user:password@localhost/rinha?sslmode=disable"
	if c := os.Getenv("DB_CONN_STR"); c != "" {
		dbConnStr = c
	}

	dbMaxOpenConns := 100
	if maxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConns != "" {
		maxOpenConnsInt, err := strconv.Atoi(maxOpenConns)
		// TODO -> log something if error is not nil
		if err == nil {
			dbMaxOpenConns = maxOpenConnsInt
		}
	}

	return config{
		dbConnStr:      dbConnStr,
		dbMaxOpenConns: dbMaxOpenConns,
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
	dbConn, err := database.NewDBConn(driverName, cfg.dbConnStr, cfg.dbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	defer dbConn.Conn.Close()

	customersRepo := repository.NewCustomers(dbConn)
	customersUsecase := customers.NewCustomerUsecase(customersRepo)

	s := server.NewServer(customersUsecase)

	s.Logger.Fatal(s.Start(bindAddr))
}
