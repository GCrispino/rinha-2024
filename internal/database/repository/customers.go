package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/GCrispino/rinha-2024/internal/database/connection"
	appErrors "github.com/GCrispino/rinha-2024/internal/errors"
	"github.com/GCrispino/rinha-2024/internal/models"
)

type Customers struct {
	dbConn *connection.DBConn
}

func NewCustomers(conn *connection.DBConn) *Customers {
	return &Customers{conn}
}

func getCustomer(ctx context.Context, tx *sql.Tx, id int) (customer *models.Customer, err error) {
	query := `SELECT id, "limit", balance, created_at from customers WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)

	customer = new(models.Customer)
	if err = row.Scan(&customer.Id, &customer.Limit, &customer.Balance, &customer.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = appErrors.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

type GetCustomerStatementResult struct {
	// customer data
	CustomerId        int
	CustomerLimit     int
	CustomerBalance   int
	CustomerCreatedAt time.Time
	// transaction data
	TransactionId          *int
	TransactionValue       *int
	TransactionType        *models.TransactionType
	TransactionDescription *string
	TransactionCustomerId  *int
	TransactionCreatedAt   *time.Time
}

func (c *Customers) GetCustomerStatement(ctx context.Context, id int) (*models.Customer, []*models.Transaction, error) {
	// get customer transactions
	query := `
		SELECT c.id, c.limit, c.balance, c.created_at, t.id, t.value, t.type, t.description, t.customer_id, t.created_at FROM customers c
		LEFT JOIN transactions t ON c.id=t.customer_id 
		WHERE c.id = $1
		ORDER BY t.created_at DESC
		LIMIT 10
	`

	rows, err := c.dbConn.Conn.QueryContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, appErrors.ErrCustomerNotFound
		}
		return nil, nil, fmt.Errorf("error getting customer transactions: %w", err)
	}
	defer rows.Close()

	results := make([]*GetCustomerStatementResult, 0)
	hasRows := false
	for rows.Next() {
		hasRows = true
		res := new(GetCustomerStatementResult)
		if err := rows.Scan(
			&res.CustomerId,
			&res.CustomerLimit,
			&res.CustomerBalance,
			&res.CustomerCreatedAt,
			&res.TransactionId,
			&res.TransactionValue,
			&res.TransactionType,
			&res.TransactionDescription,
			&res.TransactionCustomerId,
			&res.TransactionCreatedAt,
		); err != nil {
			return nil, nil, err
		}

		results = append(results, res)
	}
	if !hasRows {
		return nil, nil, appErrors.ErrCustomerNotFound
	}

	firstRes := results[0]
	customer := &models.Customer{
		Id:        firstRes.CustomerId,
		Limit:     firstRes.CustomerLimit,
		Balance:   firstRes.CustomerBalance,
		CreatedAt: firstRes.CustomerCreatedAt,
	}

	var transactions []*models.Transaction
	if len(results) == 1 && results[0].TransactionId == nil {
		return customer, transactions, nil
	}

	if firstRes.TransactionId == nil {
		return customer, transactions, nil
	}

	transactions = make([]*models.Transaction, len(results))
	for i, tx := range results {
		transactions[i] = &models.Transaction{
			Id:          *tx.TransactionId,
			Value:       *tx.TransactionValue,
			Type:        *tx.TransactionType,
			Description: *tx.TransactionDescription,
			CustomerId:  *tx.TransactionCustomerId,
			CreatedAt:   *tx.TransactionCreatedAt,
		}
	}

	return customer, transactions, nil
}

func (c *Customers) CreateCustomerTransaction(
	ctx context.Context,
	customerId int,
	value int,
	transactionType models.TransactionType,
	description string,
) (int, int, error) {
	tx, err := c.dbConn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("error beginning transaction: %w", err)
	}

	var txErr error
	defer func() {
		if txErr != nil {
			if err := tx.Rollback(); err != nil {
				// TODO -> log something
			}
		}
	}()

	updateQuery := `
		with
			c AS (SELECT * FROM customers c WHERE id = $2),
			u AS (
				UPDATE customers c2 SET balance = balance + $1
				WHERE id = $2 AND (balance + $1) >= -"limit"
				RETURNING id, "limit", balance
			),
			cu AS (SELECT COUNT(*) FROM u)
		SELECT c.limit, c.balance, cu.count as count_update FROM c, cu
    `
	insertQuery := `
      INSERT INTO transactions (value, "type", description, customer_id)
      VALUES ($1, $2, $3, $4)
    `

	updateValue := value
	if transactionType == models.TransactionTypeDebit {
		updateValue = -value
	}

	row := tx.QueryRowContext(ctx, updateQuery, updateValue, customerId)

	var limit, total, updateCount int
	if err := row.Scan(&limit, &total, &updateCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			txErr = appErrors.ErrCustomerNotFound
		} else {
			txErr = fmt.Errorf("error scanning result of customer update query: %w", err)
		}

		return 0, 0, txErr
	}

	if updateCount == 0 {
		txErr = appErrors.ErrNegativeBalanceTxResult
		return 0, 0, txErr
	}

	_, err = tx.ExecContext(ctx, insertQuery, value, transactionType, description, customerId)
	if err != nil {
		txErr = fmt.Errorf("error running transaction insert query: %w", err)
		return 0, 0, txErr
	}

	if err := tx.Commit(); err != nil {
		txErr = fmt.Errorf("error commiting transaction: %w", err)
		return 0, 0, txErr
	}

	return limit, total + updateValue, nil
}
