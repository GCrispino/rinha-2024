package customers

import (
	"context"
	"fmt"
	"time"

	"github.com/GCrispino/rinha-2024/internal/database/repository"
	"github.com/GCrispino/rinha-2024/internal/models"
)

type CustomerUsecase struct {
	// dbConn *database.DBConn
	repo *repository.Customers
}

func NewCustomerUsecase(repo *repository.Customers) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (c *CustomerUsecase) GetCustomerStatement(ctx context.Context, customerId int) (*models.GetCustomerStatementResponse, error) {
	dbCustomer, dbTxs, err := c.repo.GetCustomerStatement(ctx, customerId)
	if err != nil {
		return nil, err
	}

	txs := make([]models.StatementTransaction, len(dbTxs))
	for i := range txs {
		dbTx := dbTxs[i]
		txs[i] = models.StatementTransaction{
			Value:       dbTx.Value,
			Type:        dbTx.Type,
			Description: dbTx.Description,
			Date:        dbTx.CreatedAt,
		}
	}

	statement := models.GetCustomerStatementResponse{
		Balance: models.Balance{
			Total: dbCustomer.Balance,
			Limit: dbCustomer.Limit,
			Date:  time.Now(),
		},
		LastTransactions: txs,
	}

	return &statement, nil
}

func (c *CustomerUsecase) CreateCustomerTransaction(ctx context.Context,
	customer_id int,
	value int,
	transactionType models.TransactionType,
	description string,
) (*models.CreateCustomerTransactionResponse, error) {
	limit, total, err := c.repo.CreateCustomerTransaction(ctx, customer_id, value, transactionType, description)
	if err != nil {
		return nil, fmt.Errorf("error creating customer transaction: %w", err)
	}

	return &models.CreateCustomerTransactionResponse{
		Limit: limit,
		Total: total,
	}, nil
}
