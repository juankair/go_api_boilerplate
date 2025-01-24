package account

import (
	"context"
	"fmt"
	"time"

	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id string) (Account, error)
	Create(ctx context.Context, input FormAccountRequest) (Account, error)
	ChangePassword(ctx context.Context, id string, input ChangePasswordAccountRequest) (Account, error)
	Update(ctx context.Context, id string, input FormAccountRequest) (Account, error)
	Query(ctx context.Context, offset, limit int) ([]entity.AccountMinimalData, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id string) (Account, error)
	ToggleStatus(ctx context.Context, accountId string) (Account, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type Account struct {
	entity.Account
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type FormAccountRequest struct {
	AccountId    string `json:"account_id"`
	RoleId       string `json:"role_id"`
	EmployeeCode string `json:"employee_code"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}

type ChangePasswordAccountRequest struct {
	Password string `json:"password"`
}

func (s service) Get(ctx context.Context, id string) (Account, error) {
	account, err := s.repo.Get(ctx, id)
	if err != nil {
		return Account{}, err
	}
	return Account{account}, nil
}

func (s service) Create(ctx context.Context, req FormAccountRequest) (Account, error) {
	id := entity.GenerateID()
	now := time.Now().Format("2006-01-02 15:04:05")
	err := s.repo.Create(ctx, entity.Account{
		AccountId:    id,
		EmployeeCode: req.EmployeeCode,
		RoleId:       req.RoleId,
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		IsSuperAdmin: req.IsSuperAdmin,
		IsActive:     0,
		CreatedAt:    now,
	})
	if err != nil {
		fmt.Println(err)
		return Account{}, err
	}

	return s.Get(ctx, id)
}

func (s service) Update(ctx context.Context, accountId string, req FormAccountRequest) (Account, error) {
	account, err := s.Get(ctx, accountId)
	if err != nil {
		return account, err
	}

	account.EmployeeCode = req.EmployeeCode
	account.RoleId = req.RoleId
	account.FullName = req.FullName
	account.Email = req.Email
	account.PhoneNumber = req.PhoneNumber
	account.IsSuperAdmin = req.IsSuperAdmin

	if err := s.repo.Update(ctx, account.Account); err != nil {
		return account, err
	}

	return account, nil
}

func (s service) ChangePassword(ctx context.Context, id string, req ChangePasswordAccountRequest) (Account, error) {
	account, err := s.Get(ctx, id)
	if err != nil {
		return Account{}, err
	}

	account.Password = req.Password
	account.IsActive = 1

	return s.Get(ctx, id)
}

func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s service) Query(ctx context.Context, offset, limit int) ([]entity.AccountMinimalData, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []entity.AccountMinimalData{}
	for _, item := range items {
		result = append(result, entity.AccountMinimalData{
			AccountId:    item.AccountId,
			RoleId:       item.RoleId,
			RoleName:     item.RoleName,
			EmployeeCode: item.EmployeeCode,
			FullName:     item.FullName,
			Email:        item.Email,
			PhoneNumber:  item.PhoneNumber,
			IsSuperAdmin: item.IsSuperAdmin,
			IsActive:     item.IsActive,
			CreatedAt:    item.CreatedAt,
		})
	}
	return result, nil
}

func (s service) ToggleStatus(ctx context.Context, accountId string) (Account, error) {
	account, err := s.Get(ctx, accountId)
	if err != nil {
		return account, err
	}

	account.IsActive = 1 - account.IsActive

	if err := s.repo.Update(ctx, account.Account); err != nil {
		return account, err
	}

	return account, nil
}

func (s service) Delete(ctx context.Context, id string) (Account, error) {
	account, err := s.Get(ctx, id)
	if err != nil {
		return Account{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Account{}, err
	}
	return account, nil
}
