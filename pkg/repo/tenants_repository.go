package repo

import (
    "context"
    "fmt"
)

type Tenant struct {
    ID         int32
    Name       string
    LicenseKey string
    CreatedAt  string
    UpdatedAt  string
}

func (r *Repository) CreateTenant(ctx context.Context, name, licenseKey string) (int32, error) {
    var id int32
    err := r.db.QueryRow(ctx, "INSERT INTO tenant (name, license_key) VALUES ($1, $2) RETURNING id", name, licenseKey).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("inserting tenant: %w", err)
    }
    return id, nil
}