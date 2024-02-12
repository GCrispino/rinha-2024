package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/GCrispino/rinha/internal/usecases/customers"
)

type Server struct {
	*echo.Echo
	customers *customers.CustomerUsecase
}

func NewServer(customersUsecase *customers.CustomerUsecase) *Server {
	e := echo.New()
	e.Use(middleware.Logger())

	s := &Server{e, customersUsecase}

	s.registerHandlers()

	return s
}

func (s *Server) registerHandlers() {
	s.POST("/clientes/:id/transacoes", s.CreateTransactionHandler())
	s.GET("/clientes/:id/extrato", s.GetStatementHandler())
}
