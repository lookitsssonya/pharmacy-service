package medicine

import "errors"

var (
	ErrMedicineNotFound = errors.New("medicine not found")
	ErrNotEnoughStock   = errors.New("not enough medicine in stock")
	ErrCmdRepository    = errors.New("unable to command repository")
	ErrQueryRepository  = errors.New("unable to query repository")
)
