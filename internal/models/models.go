package models

import "time"

// DB models
type TransactionType string

const (
	TransactionTypeCredit TransactionType = "c"
	TransactionTypeDebit                  = "d"
)

type Transaction struct {
	Id          int
	Value       int
	Type        TransactionType
	Description string
	CustomerId  int
	CreatedAt   time.Time
}

type Customer struct {
	Id        int
	Limit     int
	Balance   int
	CreatedAt time.Time
}

// request/response models
type GetCustomerStatementRequest struct{}

type GetCustomerStatementResponse struct {
	Balance          Balance                `json:"saldo"`
	LastTransactions []StatementTransaction `json:"ultimas_transacoes"`
}

type Balance struct {
	Total int       `json:"total"`
	Limit int       `json:"limite"`
	Date  time.Time `json:"data_extrato"`
}

type StatementTransaction struct {
	Value       int             `json:"valor"`
	Type        TransactionType `json:"tipo"`
	Description string          `json:"descricao"`
	Date        time.Time       `json:"realizada_em"`
}

type CreateCustomerTransactionRequest struct {
	Value       int             `json:"valor"`
	Type        TransactionType `json:"tipo"`
	Description string          `json:"descricao"`
}

type CreateCustomerTransactionResponse struct {
	Limit int `json:"limite"`
	Total int `json:"saldo"`
}
