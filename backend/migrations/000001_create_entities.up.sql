CREATE TABLE entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    status VARCHAR(30) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_name_length CHECK (char_length(trim(name)) >= 3 AND char_length(name) <= 100),
    CONSTRAINT chk_description_length CHECK (char_length(description) <= 500),
    CONSTRAINT chk_type CHECK (type IN ('vehicle', 'iot_device', 'facility', 'other')),
    CONSTRAINT chk_status CHECK (status IN ('active', 'inactive', 'maintenance', 'offline')),
    CONSTRAINT chk_latitude CHECK (latitude >= -90.0 AND latitude <= 90.0),
    CONSTRAINT chk_longitude CHECK (longitude >= -180.0 AND longitude <= 180.0)
);

CREATE INDEX idx_entities_type_status ON entities (type, status);
CREATE INDEX idx_entities_location ON entities (latitude, longitude);
