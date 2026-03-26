-- ENUMS
CREATE TYPE maintenance_type AS ENUM (
    'REPAIR',
    'SERVICE',
    'INSPECTION',
    'DISPOSAL'
);

CREATE TYPE maint_status AS ENUM (
    'IN_PROGRESS',
    'COMPLETED',
    'SCRAPPED'
);

-- TABLE
CREATE TABLE asset_maintenance_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    asset_id UUID NOT NULL REFERENCES asset_inventory(id),

    maintenance_type maintenance_type NOT NULL,
    description TEXT NOT NULL,

    sent_for_repair_at DATE,
    vendor VARCHAR(200),

    maint_status maint_status DEFAULT 'IN_PROGRESS',

    returned_from_repair_at DATE,
    repair_cost_inr NUMERIC(12,2),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- INDEX
CREATE INDEX idx_asset_maintenance 
ON asset_maintenance_logs (asset_id, maint_status, created_at);
