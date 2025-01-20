package tenants

import "time"

type Tenant struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	LicenseKey string    `json:"license_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
