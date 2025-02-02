package testkit

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id int) (entity.TestKit, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, offset, limit int) ([]entity.TestKit, error)
	Create(ctx context.Context, testkit entity.TestKit) (entity.TestKit, error)
	Update(ctx context.Context, testkit entity.TestKit) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id int) (entity.TestKit, error) {
	var testkit entity.TestKit
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"id": id}).One(&testkit)
	return testkit, err
}

func (r repository) Create(ctx context.Context, testkit entity.TestKit) (entity.TestKit, error) {
	return testkit, r.db.With(ctx).Model(&testkit).Insert()
}

func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("testkit").Row(&count)
	return count, err
}

func (r repository) Update(ctx context.Context, testkit entity.TestKit) error {
	return r.db.With(ctx).Model(&testkit).Update()
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.TestKit, error) {
	var testkit []entity.TestKit
	err := r.db.With(ctx).
		Select("testkit.*").
		OrderBy("testkit.created_date").
		From("testkit").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&testkit)
	return testkit, err
}

func (r repository) Delete(ctx context.Context, id int) error {
	testkit, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&testkit).Delete()
}
