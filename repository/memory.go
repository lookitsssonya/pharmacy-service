package repository

import (
	"context"
	"sync"

	"example.com/pharmacy-service/medicine"
)

type memoryRepo struct {
	mu        sync.RWMutex
	medicines map[string]medicine.Medicine
}

func NewMemoryRepository() medicine.Repository {
	return &memoryRepo{medicines: make(map[string]medicine.Medicine)}
}

func (r *memoryRepo) CreateMedicine(ctx context.Context, m medicine.Medicine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.medicines[m.ID] = m
	return nil
}

func (r *memoryRepo) GetMedicineByID(ctx context.Context, id string) (medicine.Medicine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.medicines[id]
	if !ok {
		return medicine.Medicine{}, medicine.ErrMedicineNotFound
	}
	return m, nil
}

func (r *memoryRepo) SellMedicine(ctx context.Context, id string, count int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.medicines[id]
	if !ok {
		return medicine.ErrMedicineNotFound
	}
	if m.Quantity < count {
		return medicine.ErrNotEnoughStock
	}
	m.Quantity -= count
	r.medicines[id] = m
	return nil
}
