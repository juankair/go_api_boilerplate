package auth

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	"github.com/juankair/go_api_boilerplate/internal/entity"
	"github.com/juankair/go_api_boilerplate/pkg/dbcontext"
	"github.com/juankair/go_api_boilerplate/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, identifier string) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) error
	GetRoleName(ctx context.Context, roleId string) string
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, identifier string) (entity.Account, error) {
	var account entity.Account
	query := r.db.With(ctx).Select()

	isID := isValidID(identifier)
	if isID {
		query = query.Where(dbx.HashExp{"account_id": identifier})
	} else {
		query = query.Where(dbx.HashExp{"email": identifier})
	}

	err := query.One(&account)
	return account, err
}

func (r repository) Update(ctx context.Context, account entity.Account) error {
	return r.db.With(ctx).Model(&account).Update()
}

func (r repository) GetRoleName(ctx context.Context, roleId string) string {
	var role entity.Role
	_ = r.db.With(ctx).Select("role.*").From("role").Where(dbx.HashExp{"role_id": roleId}).One(&role)
	return role.Name
}

func isValidID(identifier string) bool {
	_, err := uuid.Parse(identifier)
	return err == nil
}
