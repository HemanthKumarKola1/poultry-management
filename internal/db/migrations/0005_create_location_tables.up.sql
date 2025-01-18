-- public.locations table
CREATE TABLE IF NOT EXISTS public.location (
    id SERIAL NOT NULL,
    tenant_id INTEGER NOT NULL REFERENCES public.tenant(id),
    location TEXT NOT NULL,
    zip_code TEXT,
    latitude DECIMAL,
    longitude DECIMAL,
    contact_person TEXT,
    phone TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, id)
);
CREATE INDEX idx_locations_tenant_id ON public.location (tenant_id);
CREATE UNIQUE INDEX idx_locations_tenant_id_location ON public.location (tenant_id, location);