package tenants

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TenantStore struct {
	db *sql.DB
}

func NewTenantStore(db *sql.DB) *TenantStore {
	return &TenantStore{db: db}
}

func (ts *TenantStore) CreateTenant(name string) (Tenant, error) {
	licenseKey := uuid.New().String()

	tx, err := ts.db.Begin() // Start a transaction
	if err != nil {
		return Tenant{}, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if any error occurs

	var tenant Tenant
	err = tx.QueryRow(`
        INSERT INTO public.tenants (name, license_key)
        VALUES ($1, $2)
        RETURNING id, created_at, updated_at
    `, name, licenseKey).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return Tenant{}, fmt.Errorf("creating tenant record: %w", err)
	}

	tenant.Name = name
	tenant.LicenseKey = licenseKey

	_, err = tx.Exec(fmt.Sprintf("CREATE SCHEMA tenant_%d;", tenant.ID))
	if err != nil {
		return Tenant{}, fmt.Errorf("creating tenant schema: %w", err)
	}

	// Create Master User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("default_master_password"), bcrypt.DefaultCost)
	if err != nil {
		return Tenant{}, fmt.Errorf("hashing master password: %w", err)
	}

	_, err = tx.Exec(`
        INSERT INTO public.users (username, password_hash, role, tenant_id)
        VALUES ($1, $2, $3, $4)
    `, "master_"+name, hashedPassword, "master", tenant.ID)
	if err != nil {
		return Tenant{}, fmt.Errorf("creating master user: %w", err)
	}

	_, err = tx.Exec(fmt.Sprintf(`
        CREATE TABLE tenant_%d.locations (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
        CREATE TABLE tenant_%d.inventory (
            id SERIAL PRIMARY KEY,
            location_id INTEGER REFERENCES tenant_%d.locations(id),
            chicken_count INTEGER,
            load_date DATE,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
        CREATE TABLE tenant_%d.feed_types (
            id SERIAL PRIMARY KEY,
            name TEXT UNIQUE NOT NULL,
            cost_per_unit DECIMAL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
        CREATE TABLE tenant_%d.feed_schedules (
            id SERIAL PRIMARY KEY,
            feed_type_id INTEGER REFERENCES tenant_%d.feed_types(id),
            times_per_day INTEGER,
            amount_per_feeding DECIMAL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
        CREATE TABLE tenant_%d.feeding_logs (
            id SERIAL PRIMARY KEY,
            feed_type_id INTEGER REFERENCES tenant_%d.feed_types(id),
            chicken_ids INTEGER[],
            amount_fed DECIMAL,
            date_time TIMESTAMP,
            comments TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
        CREATE TABLE tenant_%d.chickens (
            id SERIAL PRIMARY KEY,
            breed TEXT,
            hatch_date DATE,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );
    CREATE TABLE tenant_%d.logs (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES public.users(id),
        action TEXT,
        details JSONB,
        timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    `, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID, tenant.ID))
	if err != nil {
		return Tenant{}, fmt.Errorf("creating tenant tables: %w", err)
	}

	if err := tx.Commit(); err != nil { // Commit the transaction
		return Tenant{}, fmt.Errorf("committing transaction: %w", err)
	}

	return tenant, nil
}

// ... other methods ...
