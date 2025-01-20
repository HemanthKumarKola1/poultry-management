CREATE OR REPLACE FUNCTION create_tenant_schema(tenant_id INT)
RETURNS VOID AS $$
BEGIN
    EXECUTE format('
        CREATE SCHEMA IF NOT EXISTS tenant_%I;

        CREATE TABLE tenant_%I.locations (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.inventory (
            id SERIAL PRIMARY KEY,
            location_id INTEGER REFERENCES tenant_%I.locations(id),
            chicken_count INTEGER,
            load_date DATE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.feed_types (
            id SERIAL PRIMARY KEY,
            name TEXT UNIQUE NOT NULL,
            cost_per_unit DECIMAL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.feed_schedules (
            id SERIAL PRIMARY KEY,
            feed_type_id INTEGER REFERENCES tenant_%I.feed_types(id),
            times_per_day INTEGER,
            amount_per_feeding DECIMAL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.feeding_logs (
            id SERIAL PRIMARY KEY,
            feed_type_id INTEGER REFERENCES tenant_%I.feed_types(id),
            chicken_ids INTEGER[],
            amount_fed DECIMAL,
            date_time TIMESTAMP WITH TIME ZONE,
            comments TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.chickens (
            id SERIAL PRIMARY KEY,
            breed TEXT,
            hatch_date DATE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE tenant_%I.logs (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES public.users(id), 
            action TEXT,
            details JSONB,
            timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
    ', tenant_id); 
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION drop_tenant_schema(tenant_id INT) 
RETURNS VOID AS $$
BEGIN
    EXECUTE format('DROP SCHEMA IF EXISTS tenant_%I CASCADE', tenant_id);
END;
$$ LANGUAGE plpgsql;