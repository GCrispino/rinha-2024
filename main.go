package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	debugPprof     bool
	dbConnStr      string
	dbMaxOpenConns int
}

func loadConfig() config {
	dbConnStr := "postgres://user:password@localhost/rinha?sslmode=disable"
	if c := os.Getenv("DB_CONN_STR"); c != "" {
		dbConnStr = c
	}
	fmt.Println("dbConnStr", dbConnStr)

	debugPprof := false
	if d := os.Getenv("DEBUG_PPROF"); d != "" {
		debugPprof, _ = strconv.ParseBool(d)
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
		debugPprof:     debugPprof,
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

	ctx := context.Background()
	driverName := "postgres"
	dbConn, err := database.NewDBConn(ctx, driverName, cfg.dbConnStr, cfg.dbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	defer dbConn.Conn.Close()

	customersRepo := repository.NewCustomers(dbConn)
	customersUsecase := customers.NewCustomerUsecase(customersRepo)

	// if configured, run debug server
	if cfg.debugPprof {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// run actual server
	s := server.NewServer(customersUsecase)

	s.Logger.Fatal(s.Start(bindAddr))
}
