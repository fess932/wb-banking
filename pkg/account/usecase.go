package account

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sync"
)

func NewUsecase() *Usecase {
	return &Usecase{
		repo: make(map[string]*Account),
	}
}

type Usecase struct {
	sync.Mutex
	repo map[string]*Account
}

func (u *Usecase) Register(ctx context.Context, account *Account) error {
	// проверка на существующий
	u.Lock()
	defer u.Unlock()

	for _, v := range u.repo {
		if v.Email == account.Email {
			return ErrAccountExists
		}
	}

	// ....
	account.ID = uuid.NewString()
	u.repo[account.ID] = account

	return nil
}

func (u *Usecase) ShowAmount(ctx context.Context, accountID string) (decimal.Decimal, error) {
	u.Lock()
	defer u.Unlock()

	v, ok := u.repo[accountID]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("account id %v: %w", accountID, ErrAccountNotExists)
	}

	return v.Amount, nil
}

func (u *Usecase) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount decimal.Decimal) error {
	u.Lock()
	defer u.Unlock()

	from, ok := u.repo[fromAccountID]
	if !ok {
		return fmt.Errorf("account id %v: %w", fromAccountID, ErrAccountNotExists)
	}

	to, ok := u.repo[toAccountID]
	if !ok {
		return fmt.Errorf("account id %v: %w", toAccountID, ErrAccountNotExists)
	}

	if err := from.Transfer(to, amount); err != nil {
		return fmt.Errorf("cant transfer: %w", err)
	}

	return nil
}
