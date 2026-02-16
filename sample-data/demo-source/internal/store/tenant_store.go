package store

import "errors"

type TenantRecord struct {
	ID     int
	Name   string
	Active bool
}

func FindTenantByID(id int) (TenantRecord, error) {
	if id <= 0 {
		return TenantRecord{}, errors.New("invalid tenant id")
	}
	return TenantRecord{ID: id, Name: "Acme", Active: true}, nil
}
