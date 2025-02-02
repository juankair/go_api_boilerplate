package testkit

import (
	"context"
	"fmt"
	"time"

	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id int) (TestKit, error)
	Create(ctx context.Context, input FormTestKitRequest) (TestKit, error)
	Update(ctx context.Context, id int, input FormTestKitRequest) (TestKit, error)
	Query(ctx context.Context, offset, limit int) ([]entity.TestKit, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int) (TestKit, error)
	ToggleStatus(ctx context.Context, id int) (TestKit, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type TestKit struct {
	entity.TestKit
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type FormTestKitRequest struct {
	ID      int    `json:"id"`
	TestKit string `json:"testkit"`
	Hasil   string `json:"hasil"`
}

func (s service) Get(ctx context.Context, id int) (TestKit, error) {
	testkit, err := s.repo.Get(ctx, id)
	if err != nil {
		return TestKit{}, err
	}
	return TestKit{testkit}, nil
}

func (s service) Create(ctx context.Context, req FormTestKitRequest) (TestKit, error) {
	accountID, _ := ctx.Value("accountID").(string)
	now := time.Now().Format("2006-01-02 15:04:05")
	res, err := s.repo.Create(ctx, entity.TestKit{
		Testkit:     req.TestKit,
		Hasil:       req.Hasil,
		CreatedBy:   &accountID,
		CreatedDate: &now,
	})
	if err != nil {
		fmt.Println(err)
		return TestKit{}, err
	}

	return s.Get(ctx, res.ID)
}

func (s service) Update(ctx context.Context, id int, req FormTestKitRequest) (TestKit, error) {
	testkit, err := s.Get(ctx, id)
	if err != nil {
		return testkit, err
	}

	testkit.Testkit = req.TestKit
	testkit.Hasil = req.Hasil

	if err := s.repo.Update(ctx, testkit.TestKit); err != nil {
		return testkit, err
	}

	return testkit, nil
}

func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s service) Query(ctx context.Context, offset, limit int) ([]entity.TestKit, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []entity.TestKit{}
	for _, item := range items {
		result = append(result, entity.TestKit{
			ID:          item.ID,
			Testkit:     item.Testkit,
			Hasil:       item.Hasil,
			Deleted:     item.Deleted,
			CreatedBy:   item.CreatedBy,
			CreatedDate: item.CreatedDate,
			UpdatedBy:   item.UpdatedBy,
			UpdatedDate: item.UpdatedDate,
		})
	}
	return result, nil
}

func (s service) ToggleStatus(ctx context.Context, id int) (TestKit, error) {
	testkit, err := s.Get(ctx, id)
	if err != nil {
		return testkit, err
	}

	testkit.Deleted = 1 - testkit.Deleted

	if err := s.repo.Update(ctx, testkit.TestKit); err != nil {
		return testkit, err
	}

	return testkit, nil
}

func (s service) Delete(ctx context.Context, id int) (TestKit, error) {
	testkit, err := s.Get(ctx, id)
	if err != nil {
		return TestKit{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return TestKit{}, err
	}
	return testkit, nil
}
