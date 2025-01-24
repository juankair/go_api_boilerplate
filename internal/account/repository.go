package account

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Account, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]entity.AccountMinimalData, error)
	Create(ctx context.Context, account entity.Account) error
	Update(ctx context.Context, account entity.Account) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, accountId string) (entity.Account, error) {
	var account entity.Account
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"account_id": accountId}).One(&account)
	return account, err
}

func (r repository) Create(ctx context.Context, account entity.Account) error {
	return r.db.With(ctx).Model(&account).Insert()
}

func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("account").Row(&count)
	return count, err
}

func (r repository) Update(ctx context.Context, account entity.Account) error {
	return r.db.With(ctx).Model(&account).Update()
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.AccountMinimalData, error) {
	var account []entity.AccountMinimalData
	err := r.db.With(ctx).
		Select("account.*", "COALESCE(role.name,'') as role_name").
		Join("LEFT JOIN", "role", dbx.NewExp("account.role_id = role.role_id")).
		OrderBy("account.created_at").
		From("account").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&account)
	return account, err
}

func (r repository) Delete(ctx context.Context, id string) error {
	account, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&account).Delete()
}
