-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_logs (
    id SERIAL PRIMARY KEY,
    product_id UUID NOT NULL,
    old_name TEXT,
    new_name TEXT,
    old_description TEXT,
    new_description TEXT,
    old_quantity INTEGER,
    new_quantity INTEGER,
    changed_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_logs;
-- +goose StatementEnd
