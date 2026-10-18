package medicine

import "context"

type Medicine struct {
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name"`
	Manufacturer string  `json:"manufacturer"`
	Price        float64 `json:"price"`
	Prescription bool    `json:"prescription"`
	Quantity     int     `json:"quantity"`
}

type Service interface {
	Create(ctx context.Context, medicine Medicine) (string, error)
	GetByID(ctx context.Context, id string) (Medicine, error)
	Sell(ctx context.Context, id string, count int) error
}

type Repository interface {
	CreateMedicine(ctx context.Context, m Medicine) error
	GetMedicineByID(ctx context.Context, id string) (Medicine, error)
	SellMedicine(ctx context.Context, id string, count int) error
}
