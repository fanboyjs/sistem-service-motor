CREATE TABLE service_records (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    service_type_id BIGINT NOT NULL REFERENCES service_types(id) ON DELETE CASCADE,
    service_date DATE NOT NULL,
    odometer DECIMAL(15, 2) NOT NULL,
    workshop_name VARCHAR(150) NOT NULL,
    labor_cost DECIMAL(15, 2) NOT NULL,
    parts_cost DECIMAL(15, 2) NOT NULL,
    total_cost DECIMAL(15, 2) NOT NULL,
    notes TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
