CREATE TABLE IF NOT EXISTS estimate_items (
    line_item_id SERIAL PRIMARY KEY,
    estimate_id INT NOT NULL REFERENCES estimates(estimate_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id),
    quantity INT NOT NULL DEFAULT 1
);
