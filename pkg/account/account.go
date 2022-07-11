package account

import (
	"errors"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID     string
	Email  string
	Amount decimal.Decimal
}

func (a *Account) Transfer(to *Account, amount decimal.Decimal) error {
	n := a.Amount.Sub(amount)
	if n.IsNegative() {
		return ErrInsufficientFunds
	}

	a.Amount = n
	to.Amount = to.Amount.Add(amount)

	return nil
}

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrAccountExists = errors.New("account exists")
var ErrAccountNotExists = errors.New("account not exists")
