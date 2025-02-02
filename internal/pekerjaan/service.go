package pekerjaan

import (
	"context"
	"fmt"
	"time"

	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id int) (Pekerjaan, error)
	Create(ctx context.Context, input FormPekerjaanRequest) (Pekerjaan, error)
	Update(ctx context.Context, id int, input FormPekerjaanRequest) (Pekerjaan, error)
	Query(ctx context.Context, offset, limit int) ([]entity.Pekerjaan, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int) (Pekerjaan, error)
	ToggleStatus(ctx context.Context, id int) (Pekerjaan, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type Pekerjaan struct {
	entity.Pekerjaan
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type FormPekerjaanRequest struct {
	ID        int    `json:"id"`
	Pekerjaan string `json:"pekerjaan"`
	Note      string `json:"note"`
}

func (s service) Get(ctx context.Context, id int) (Pekerjaan, error) {
	pekerjaan, err := s.repo.Get(ctx, id)
	if err != nil {
		return Pekerjaan{}, err
	}
	return Pekerjaan{pekerjaan}, nil
}

func (s service) Create(ctx context.Context, req FormPekerjaanRequest) (Pekerjaan, error) {
	accountID, _ := ctx.Value("accountID").(string)
	now := time.Now().Format("2006-01-02 15:04:05")
	res, err := s.repo.Create(ctx, entity.Pekerjaan{
		Pekerjaan:   req.Pekerjaan,
		Note:        &req.Note,
		CreatedBy:   &accountID,
		CreatedDate: &now,
	})
	if err != nil {
		fmt.Println(err)
		return Pekerjaan{}, err
	}

	return s.Get(ctx, res.ID)
}

func (s service) Update(ctx context.Context, id int, req FormPekerjaanRequest) (Pekerjaan, error) {
	pekerjaan, err := s.Get(ctx, id)
	if err != nil {
		return pekerjaan, err
	}

	pekerjaan.Pekerjaan.Pekerjaan = req.Pekerjaan
	pekerjaan.Note = &req.Note

	if err := s.repo.Update(ctx, pekerjaan.Pekerjaan); err != nil {
		return pekerjaan, err
	}

	return pekerjaan, nil
}

func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s service) Query(ctx context.Context, offset, limit int) ([]entity.Pekerjaan, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []entity.Pekerjaan{}
	for _, item := range items {
		result = append(result, entity.Pekerjaan{
			ID:          item.ID,
			Pekerjaan:   item.Pekerjaan,
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

func (s service) ToggleStatus(ctx context.Context, id int) (Pekerjaan, error) {
	pekerjaan, err := s.Get(ctx, id)
	if err != nil {
		return pekerjaan, err
	}

	pekerjaan.Deleted = 1 - pekerjaan.Deleted

	if err := s.repo.Update(ctx, pekerjaan.Pekerjaan); err != nil {
		return pekerjaan, err
	}

	return pekerjaan, nil
}

func (s service) Delete(ctx context.Context, id int) (Pekerjaan, error) {
	pekerjaan, err := s.Get(ctx, id)
	if err != nil {
		return Pekerjaan{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Pekerjaan{}, err
	}
	return pekerjaan, nil
}
