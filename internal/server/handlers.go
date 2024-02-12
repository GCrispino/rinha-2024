package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/GCrispino/rinha/internal/models"
	"github.com/labstack/echo/v4"
)

func getCustomerIdFromRequest(c echo.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// TODO -> return error
		return 0, err
	}
	return id, nil
}

func (s *Server) CreateTransactionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		customerId, err := getCustomerIdFromRequest(c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		req := new(models.CreateCustomerTransactionRequest)
		if err := c.Bind(req); err != nil {
			return fmt.Errorf("error binding CreateCustomerTransactionRequest request: %w", err)
		}

		res, err := s.customers.CreateCustomerTransaction(
			ctx,
			customerId,
			req.Value, req.Type, req.Description,
		)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (s *Server) GetStatementHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := getCustomerIdFromRequest(c)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		resp, err := s.customers.GetCustomerStatement(ctx, id)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}
}
