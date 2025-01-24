package document

import (
	"context"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) (res *entity.Document, err error)
	//Count(ctx context.Context) (int, error)
	//Query(ctx context.Context, offset, limit int) ([]entity.Document, error)
	//Create(ctx context.Context, account entity.Document) error
	//Update(ctx context.Context, account entity.Document) error
	//Delete(ctx context.Context, id string) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id string) (res *entity.Document, err error) {
	return res, err
}
