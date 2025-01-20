package tenants

import "database/sql"

type Service struct {
	store *TenantStore
}

func NewService(db *sql.DB) *Service {
	return &Service{store: NewTenantStore(db)}
}

func (s *Service) CreateTenant(name string) (Tenant, error) {
	return s.store.CreateTenant(name)
}
