package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	appErrors "github.com/GCrispino/rinha-2024/internal/errors"
	"github.com/GCrispino/rinha-2024/internal/models"
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
			var unmarshalErr *json.UnmarshalTypeError
			if errors.As(err, &unmarshalErr) {
				return echo.ErrUnprocessableEntity
			}
			return fmt.Errorf("error binding CreateCustomerTransactionRequest request: %w", err)
		}

		res, err := s.customers.CreateCustomerTransaction(
			ctx,
			customerId,
			req.Value, req.Type, req.Description,
		)
		if err != nil {
			switch {
			case errors.Is(err, appErrors.ErrNegativeBalanceTxResult):
				return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
			case errors.Is(err, appErrors.ErrCustomerNotFound):
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			default:
				return err
			}
		}

		s.Echo.Logger.Debug("response:", res)
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
			if errors.Is(err, appErrors.ErrCustomerNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}
}
