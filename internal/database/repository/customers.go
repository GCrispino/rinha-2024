package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GCrispino/rinha-2024/internal/database/connection"
	"github.com/GCrispino/rinha-2024/internal/models"
)

type Customers struct {
	dbConn *connection.DBConn
}

func NewCustomers(conn *connection.DBConn) *Customers {
	return &Customers{conn}
}

func (c *Customers) GetCustomerStatement(ctx context.Context, id int) (*models.Customer, []*models.Transaction, error) {
	tx, err := c.dbConn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	var txErr error
	defer func() {
		if txErr != nil {
			if err := tx.Rollback(); err != nil {
				// TODO -> log something
			}
		}
	}()

	query := `SELECT id, "limit", balance, created_at from customers WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)

	var customer models.Customer
	if err := row.Scan(&customer.Id, &customer.Limit, &customer.Balance, &customer.CreatedAt); err != nil {
		txErr = err
		return nil, nil, txErr
	}

	// get customer transactions
	query = `SELECT id, value, type, customer_id, created_at from transactions WHERE customer_id = $1 LIMIT 10`
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		txErr = fmt.Errorf("error getting customer transactions: %w", err)
		return nil, nil, txErr
	}

	transactions := make([]*models.Transaction, 0)
	for rows.Next() {
		transaction := new(models.Transaction)
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Value,
			&transaction.Type,
			&transaction.CustomerId,
			&transaction.CreatedAt,
		); err != nil {
			return nil, nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := tx.Commit(); err != nil {
		txErr = fmt.Errorf("error commiting transaction: %w", err)
		return nil, nil, txErr
	}

	return &customer, transactions, nil
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
      UPDATE customers SET balance = balance + $1
      WHERE id = $2 AND (balance + $1) >= -"limit"
    `
	insertQuery := `
      INSERT INTO transactions (value, "type", customer_id)
      VALUES ($1, $2, $3)
    `

	updateValue := value
	if transactionType == models.TransactionTypeDebit {
		updateValue = -value
	}

	// res, err := tx.ExecContext(ctx, updateQuery, updateValue, customerId)
	// if err != nil {
	// 	txErr = fmt.Errorf("error running customer updating query: %w", err)
	// 	return 0, 0, txErr
	// }

	// rowsAffected, err := res.RowsAffected()
	// if err != nil {
	// 	txErr = fmt.Errorf("error obtaining rows affected by customer updating query: %w", err)
	// 	return 0, 0, txErr
	// }

	// if rowsAffected == 0 {
	// 	txErr = fmt.Errorf("transaction results in negative balance")
	// 	return 0, 0, txErr
	// }

	row := tx.QueryRowContext(ctx, updateQuery, updateValue, customerId)

	var limit, total int
	if err := row.Scan(&limit, &total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			txErr = fmt.Errorf("transaction results in negative balance")
		} else {
			txErr = fmt.Errorf("error scanning result of customer update query: %w", err)
		}

		return 0, 0, txErr
	}

	_, err = tx.ExecContext(ctx, insertQuery, value, transactionType, customerId)
	if err != nil {
		txErr = fmt.Errorf("error running transaction insert query: %w", err)
		return 0, 0, txErr
	}

	if err := tx.Commit(); err != nil {
		txErr = fmt.Errorf("error commiting transaction: %w", err)
		return 0, 0, txErr
	}

	return limit, total, nil
}
