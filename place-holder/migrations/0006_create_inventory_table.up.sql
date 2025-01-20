CREATE TABLE IF NOT EXISTS public.inventory (
    id SERIAL NOT NULL,
    tenant_id INTEGER NOT NULL REFERENCES public.tenant(id),
    location_id INTEGER NOT NULL,
    chicken_count INTEGER,
    feed INTEGER,
    load_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, id),
    FOREIGN KEY (location_id, tenant_id) REFERENCES public.location(id, tenant_id),
    UNIQUE (location_id, tenant_id) -- Ensures 1:1 relationship
);
CREATE INDEX idx_inventory_tenant_id ON public.inventory (tenant_id);
CREATE INDEX idx_inventory_location_id_tenant_id ON public.inventory (location_id, tenant_id);