-- ENUMS
CREATE TYPE asset_type AS ENUM ('HARDWARE','SOFTWARE','FURNITURE','OTHER');

CREATE TYPE asset_status AS ENUM (
  'AVAILABLE','ASSIGNED','UNDER_REPAIR','RETIRED','LOST'
);

CREATE TYPE asset_category AS ENUM (
  'LAPTOP','MOBILE','DESKTOP','CHAIR','TABLE','OTHER'
);

-- SEQUENCE
CREATE SEQUENCE asset_code_seq START 1;

-- TABLE
CREATE TABLE asset_inventory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  asset_code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  serial_no TEXT,
  asset_type asset_type NOT NULL,
  category asset_category,
  status asset_status DEFAULT 'AVAILABLE',
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- UNIQUE NULLS NOT DISTINCT
CREATE UNIQUE INDEX ux_asset_serial_no
ON asset_inventory(serial_no)
WHERE serial_no IS NOT NULL;

-- INDEX
CREATE INDEX idx_asset_filters
ON asset_inventory(status, asset_type, is_deleted);