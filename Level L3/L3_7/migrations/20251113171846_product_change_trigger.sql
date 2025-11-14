-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_product_update()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO product_logs (
    product_id,
    old_name, new_name,
    old_description, new_description,
    old_quantity, new_quantity,
    changed_at
)
VALUES (
           OLD.product_id,
           OLD.product_name, NEW.product_name,
           OLD.description, NEW.description,
           OLD.quantity, NEW.quantity,
           NOW()
       );
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_product_update_trigger
    AFTER UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION log_product_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS log_product_update_trigger ON products;
DROP FUNCTION IF EXISTS log_product_update;
-- +goose StatementEnd
