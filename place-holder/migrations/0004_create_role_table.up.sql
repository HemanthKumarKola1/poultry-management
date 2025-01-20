CREATE TABLE public.role (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    permissions TEXT[] NOT NULL
);

INSERT INTO public.role (name, permissions) VALUES
('master', '{users:create,users:read,users:update,users:delete,inventory:create,inventory:read,inventory:update,inventory:delete,feed:create,feed:read,feed:update,feed:delete,analytics:read}'),
('admin', '{users:create,users:read,inventory:create,inventory:read,inventory:update,feed:create,feed:read,feed:update,analytics:read}'),
('farm_worker', '{feed:create,feed:read}'),
('viewer', '{inventory:read,feed:read,analytics:read}');