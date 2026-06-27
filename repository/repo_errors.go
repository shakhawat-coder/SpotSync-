package repository

import "errors"

var (
	ErrZoneAtCapacity         = errors.New("zone is at full capacity")
	ErrDuplicateActiveLicense = errors.New("duplicate active license plate")
)
