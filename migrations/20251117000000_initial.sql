-- +goose Up
-- +goose StatementBegin
-- Current configuration (single row, always id=1)
CREATE TABLE pack_configuration (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    version INTEGER NOT NULL DEFAULT 1,
    pack_sizes JSONB NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255)
);

-- Historical configurations (audit trail)
CREATE TABLE pack_configuration_history (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL,
    pack_sizes JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)
);

-- Initialize with default config
INSERT INTO pack_configuration (id, pack_sizes, updated_by) 
VALUES (1, '[250, 500, 1000, 2000, 5000]', 'system');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pack_configuration_history;
DROP TABLE IF EXISTS pack_configuration;
-- +goose StatementEnd
