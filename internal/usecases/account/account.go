package account

import (
	"errors"
	"fmt"
	"github.com/mp-hl-2021/muzio/internal/domain/account"
	"github.com/mp-hl-2021/muzio/internal/service/auth"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

const (
	minLoginLength    = 3
	maxLoginLength    = 20
	minPasswordLength = 3
	maxPasswordLength = 48
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
	AuthToken 	   auth.Interface
}

func (a *UseCases) CreateAccount(login, password string) (Account, error) {
	err := validateLogin(login)
	if err != nil {
		return Account{}, err
	}
	err = validatePassword(password)
	if err != nil {
		return Account{}, err
	}
	hasedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	acc, err := a.AccountStorage.CreateAccount(account.Credentials{
		Login: login,
		Password: string(hasedPass),
	})
	if err != nil {
		return Account{}, err
	}

	return Account{Id: acc.Id}, nil
}

func (a *UseCases) GetAccountById(id string) (Account, error) {
	acc, err := a.AccountStorage.GetAccountById(id)
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, nil
}

func (a *UseCases) LoginToAccount(login, password string) (string, error) {
	err := validateLogin(login)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	err = validatePassword(password)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	acc, err := a.AccountStorage.GetAccountByLogin(login)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password))
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	t, err := a.AuthToken.IssueToken(acc.Id)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (a *UseCases) Authenticate(token string) (string, error) {
	return a.AuthToken.UserIdByToken(token)
}

func validateLogin (login string) error {
	return validate(login, "login", maxLoginLength, minLoginLength)
}

func validatePassword (pass string) error {
	return validate(pass, "password", maxPasswordLength, minPasswordLength)
}

func validate (s, strType string, maxLen int, minLen int) error {
	i := 0

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.New(fmt.Sprintf("Invalid %s symbol", strType))
		}
		i++
	}

	if i < minLen || i > maxLen {
		return errors.New(fmt.Sprintf("Invalid %s length", strType))
	}
	return nil
}