CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    service_date DATE NOT NULL,
    mileage INTEGER NOT NULL,
    complain TEXT NOT NULL,
    diagnosis VARCHAR(100) NOT NULL,
    notes TEXT NOT NULL,
    next_service_date DATE NOT NULL,
    next_service_mileage INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);