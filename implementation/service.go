package implementation

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"example.com/pharmacy-service/medicine"
)

type service struct {
	repository medicine.Repository
	logger     log.Logger
}

func NewService(repo medicine.Repository, logger log.Logger) medicine.Service {
	return &service{repository: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, m medicine.Medicine) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	m.ID = id

	if err := s.repository.CreateMedicine(ctx, m); err != nil {
		level.Error(logger).Log("err", err)
		return "", medicine.ErrCmdRepository
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (medicine.Medicine, error) {
	logger := log.With(s.logger, "method", "GetByID")

	m, err := s.repository.GetMedicineByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return m, medicine.ErrMedicineNotFound
	}
	return m, nil
}

func (s *service) Sell(ctx context.Context, id string, count int) error {
	logger := log.With(s.logger, "method", "Sell")

	if count <= 0 {
		return medicine.ErrNotEnoughStock
	}

	if err := s.repository.SellMedicine(ctx, id, count); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}
	return nil
}
