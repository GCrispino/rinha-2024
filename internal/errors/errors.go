package errors

import "fmt"

var ErrNegativeBalanceTxResult = fmt.Errorf("transaction results in negative balance")
var ErrCustomerNotFound = fmt.Errorf("customer not found")
