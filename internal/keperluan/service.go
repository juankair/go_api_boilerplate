package keperluan

import (
	"context"
	"fmt"
	"time"

	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id int) (Keperluan, error)
	Create(ctx context.Context, input FormKeperluanRequest) (Keperluan, error)
	Update(ctx context.Context, id int, input FormKeperluanRequest) (Keperluan, error)
	Query(ctx context.Context, offset, limit int) ([]entity.Keperluan, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int) (Keperluan, error)
	ToggleStatus(ctx context.Context, id int) (Keperluan, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type Keperluan struct {
	entity.Keperluan
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type FormKeperluanRequest struct {
	ID        int    `json:"id"`
	Keperluan string `json:"keperluan"`
	Note      string `json:"note"`
}

func (s service) Get(ctx context.Context, id int) (Keperluan, error) {
	keperluan, err := s.repo.Get(ctx, id)
	if err != nil {
		return Keperluan{}, err
	}
	return Keperluan{keperluan}, nil
}

func (s service) Create(ctx context.Context, req FormKeperluanRequest) (Keperluan, error) {
	accountID, _ := ctx.Value("accountID").(string)
	now := time.Now().Format("2006-01-02 15:04:05")
	res, err := s.repo.Create(ctx, entity.Keperluan{
		Keperluan:   req.Keperluan,
		Note:        &req.Note,
		CreatedBy:   &accountID,
		CreatedDate: &now,
	})
	if err != nil {
		fmt.Println(err)
		return Keperluan{}, err
	}

	return s.Get(ctx, res.ID)
}

func (s service) Update(ctx context.Context, id int, req FormKeperluanRequest) (Keperluan, error) {
	keperluan, err := s.Get(ctx, id)
	if err != nil {
		return keperluan, err
	}

	keperluan.Keperluan.Keperluan = req.Keperluan
	keperluan.Note = &req.Note

	if err := s.repo.Update(ctx, keperluan.Keperluan); err != nil {
		return keperluan, err
	}

	return keperluan, nil
}

func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s service) Query(ctx context.Context, offset, limit int) ([]entity.Keperluan, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []entity.Keperluan{}
	for _, item := range items {
		result = append(result, entity.Keperluan{
			ID:          item.ID,
			Keperluan:   item.Keperluan,
			Note:        item.Note,
			Deleted:     item.Deleted,
			CreatedBy:   item.CreatedBy,
			CreatedDate: item.CreatedDate,
			UpdatedBy:   item.UpdatedBy,
			UpdatedDate: item.UpdatedDate,
		})
	}
	return result, nil
}

func (s service) ToggleStatus(ctx context.Context, id int) (Keperluan, error) {
	keperluan, err := s.Get(ctx, id)
	if err != nil {
		return keperluan, err
	}

	keperluan.Deleted = 1 - keperluan.Deleted

	if err := s.repo.Update(ctx, keperluan.Keperluan); err != nil {
		return keperluan, err
	}

	return keperluan, nil
}

func (s service) Delete(ctx context.Context, id int) (Keperluan, error) {
	keperluan, err := s.Get(ctx, id)
	if err != nil {
		return Keperluan{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Keperluan{}, err
	}
	return keperluan, nil
}
