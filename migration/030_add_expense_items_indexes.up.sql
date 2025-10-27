CREATE INDEX idx_expense_items_on_expense_id ON expense_items (expense_id);
CREATE INDEX idx_expense_items_on_product_id ON expense_items (product_id);
CREATE INDEX idx_expense_items_on_expense_id_and_is_arrived ON expense_items (expense_id, is_arrived);