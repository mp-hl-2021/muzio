package account

import (
	"github.com/mp-hl-2021/muzio/internal/domain/account"
)

type Account struct {
	Id string
}

type Interface interface {
	CreateAccount(login, password string) (Account, error)
	GetAccountById(id string) (Account, error)

	LoginToAccount(login, password string) (string, error)
	Authenticate(token string) (string, error)
}

type UseCases struct {
	AccountStorage account.Interface
	// TODO Token
}

func (a *UseCases) CreateAccount(login, password string) (Account, error) {
	panic("implement me")
}

func (a *UseCases) GetAccountById(id string) (Account, error) {
	panic("implement me")
}

func (a *UseCases) LoginToAccount(login, password string) (string, error) {
	panic("implement me")
}

func (a *UseCases) Authenticate(token string) (string, error) {
	panic("implement me")
}
