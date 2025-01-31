package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Get(ctx context.Context, email string) (entity.Account, error)
	Login(ctx context.Context, email, password string) (entity.AccountMinimalData, error)
	CheckActivation(ctx context.Context, email string) (entity.AccountMinimalData, error)
	ActivationConfirmation(ctx context.Context, email string, password string) (entity.AccountMinimalData, error)
}

type Identity interface {
	GetID() string
	GetName() string
}

type service struct {
	repo            Repository
	signingKey      string
	tokenExpiration int
	logger          log.Logger
}

type ActivationConfirmationRequest struct {
	Password string `json:"password"`
}

func NewService(repo Repository, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{repo, signingKey, tokenExpiration, logger}
}

func (s service) Get(ctx context.Context, email string) (entity.Account, error) {
	account, err := s.repo.Get(ctx, email)
	if err != nil {
		return entity.Account{}, err
	}

	return account, nil
}

func (s service) Login(ctx context.Context, email, password string) (entity.AccountMinimalData, error) {
	identity, err := s.authenticate(ctx, email, password)
	if err != nil {
		return entity.AccountMinimalData{}, err
	}
	return identity, nil
}

func (s service) CheckActivation(ctx context.Context, email string) (entity.AccountMinimalData, error) {
	account, err := s.repo.Get(ctx, email)
	if err != nil {
		return entity.AccountMinimalData{}, errors.New("Akun tidak ditemukan")
	}

	if account.IsActive == 1 {
		return entity.AccountMinimalData{}, errors.New("Akun sudah dikonfirmasi")
	}

	return entity.AccountMinimalData{
		AccountId:    account.AccountId,
		EmployeeCode: account.EmployeeCode,
		RoleName:     s.repo.GetRoleName(ctx, account.RoleId),
		FullName:     account.FullName,
		Email:        account.Email,
		PhoneNumber:  account.PhoneNumber,
		IsSuperAdmin: account.IsSuperAdmin,
	}, nil
}

func (s service) ActivationConfirmation(ctx context.Context, email string, password string) (entity.AccountMinimalData, error) {
	account, err := s.Get(ctx, email)
	if err != nil {
		return entity.AccountMinimalData{}, err
	}

	account.Password = password
	account.IsActive = 1
	token, _ := s.generateJWT(account)

	if err := s.repo.Update(ctx, account); err != nil {
		fmt.Println(err)
		return entity.AccountMinimalData{}, err
	}

	return entity.AccountMinimalData{
		AccountId:    account.AccountId,
		EmployeeCode: account.EmployeeCode,
		RoleName:     s.repo.GetRoleName(ctx, account.RoleId),
		FullName:     account.FullName,
		Email:        account.Email,
		PhoneNumber:  account.PhoneNumber,
		Token:        token,
		IsSuperAdmin: account.IsSuperAdmin,
	}, nil
}

func (s service) authenticate(ctx context.Context, email, password string) (entity.AccountMinimalData, error) {
	account, err := s.Get(ctx, email)
	if err != nil {
		return entity.AccountMinimalData{}, errors.New("Akun tidak ditemukan")
	}

	if account.Password == "" {
		return entity.AccountMinimalData{}, errors.New("Akun belum dikonfirmasi")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return entity.AccountMinimalData{}, errors.New("Password Salah")
	}

	token, _ := s.generateJWT(account)
	return entity.AccountMinimalData{
		AccountId:    account.AccountId,
		EmployeeCode: account.EmployeeCode,
		RoleName:     s.repo.GetRoleName(ctx, account.RoleId),
		FullName:     account.FullName,
		Email:        account.Email,
		PhoneNumber:  account.PhoneNumber,
		Token:        token,
		IsSuperAdmin: account.IsSuperAdmin,
	}, nil

}

func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id:        identity.GetID(),
		Audience:  identity.GetName(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
