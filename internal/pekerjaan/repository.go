package pekerjaan

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id int) (entity.Pekerjaan, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]entity.Pekerjaan, error)
	Create(ctx context.Context, pekerjaan entity.Pekerjaan) (entity.Pekerjaan, error)
	Update(ctx context.Context, pekerjaan entity.Pekerjaan) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id int) (entity.Pekerjaan, error) {
	var pekerjaan entity.Pekerjaan
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"id": id}).One(&pekerjaan)
	return pekerjaan, err
}

func (r repository) Create(ctx context.Context, pekerjaan entity.Pekerjaan) (entity.Pekerjaan, error) {
	return pekerjaan, r.db.With(ctx).Model(&pekerjaan).Insert()
}

func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("pekerjaan").Row(&count)
	return count, err
}

func (r repository) Update(ctx context.Context, pekerjaan entity.Pekerjaan) error {
	return r.db.With(ctx).Model(&pekerjaan).Update()
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Pekerjaan, error) {
	var pekerjaan []entity.Pekerjaan
	err := r.db.With(ctx).
		Select("pekerjaan.*").
		OrderBy("pekerjaan.created_date").
		From("pekerjaan").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&pekerjaan)
	return pekerjaan, err
}

func (r repository) Delete(ctx context.Context, id int) error {
	pekerjaan, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&pekerjaan).Delete()
}
