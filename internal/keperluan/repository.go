package keperluan

import (
	"context"
)

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id int) (entity.Keperluan, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]entity.Keperluan, error)
	Create(ctx context.Context, keperluan entity.Keperluan) (entity.Keperluan, error)
	Update(ctx context.Context, keperluan entity.Keperluan) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id int) (entity.Keperluan, error) {
	var keperluan entity.Keperluan
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"id": id}).One(&keperluan)
	return keperluan, err
}

func (r repository) Create(ctx context.Context, keperluan entity.Keperluan) (entity.Keperluan, error) {
	return keperluan, r.db.With(ctx).Model(&keperluan).Insert()
}

func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("keperluan").Row(&count)
	return count, err
}

func (r repository) Update(ctx context.Context, keperluan entity.Keperluan) error {
	return r.db.With(ctx).Model(&keperluan).Update()
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Keperluan, error) {
	var keperluan []entity.Keperluan
	err := r.db.With(ctx).
		Select("keperluan.*").
		//Join("LEFT JOIN", "role", dbx.NewExp("keperluan.role_id = role.role_id")).
		OrderBy("keperluan.created_date").
		From("keperluan").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&keperluan)
	return keperluan, err
}

func (r repository) Delete(ctx context.Context, id int) error {
	keperluan, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&keperluan).Delete()
}
