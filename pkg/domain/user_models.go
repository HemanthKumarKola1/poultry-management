package domain

import "github.com/jackc/pgx/v5/pgtype"

type User struct { //Includes all types of users
	ID        int32              `json:"id"`
	Username  string             `json:"username"`
	Role      string             `json:"role"`
	TenantID  int32              `json:"tenant_id"` // should be zero atleast [super_admin usecase]
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
