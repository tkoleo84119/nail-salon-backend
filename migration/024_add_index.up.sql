CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_bookings_on_customer_id ON bookings (customer_id);
CREATE INDEX idx_bookings_on_time_slot_id ON bookings (time_slot_id);
CREATE INDEX idx_bookings_on_stylist_id ON bookings (stylist_id);
CREATE INDEX idx_bookings_on_store_id_and_status ON bookings (store_id, status);
CREATE INDEX idx_bookings_comprehensive ON bookings (store_id, stylist_id, status, created_at);
CREATE INDEX idx_booking_details_on_booking_id ON booking_details (booking_id);

CREATE INDEX idx_schedules_on_work_date ON schedules (work_date);
CREATE INDEX idx_schedules_on_store_id_and_stylist_id ON schedules (store_id, stylist_id);

CREATE INDEX idx_time_slots_on_schedule_id ON time_slots (schedule_id);
CREATE INDEX idx_time_slots_on_schedule_id_and_is_available ON time_slots (schedule_id, is_available);

CREATE INDEX idx_account_transactions_on_account_id ON account_transactions (account_id);
CREATE INDEX idx_checkouts_on_booking_id ON checkouts (booking_id);

CREATE INDEX idx_staff_user_store_access_on_store_id ON staff_user_store_access (store_id);
CREATE INDEX idx_staff_user_store_access_on_staff_user_id ON staff_user_store_access (staff_user_id);

CREATE INDEX idx_customers_on_name_trigram ON customers USING gin (name gin_trgm_ops);
CREATE INDEX idx_customers_on_phone_trigram ON customers USING gin (phone gin_trgm_ops);
CREATE INDEX idx_customers_on_line_name_trigram ON customers USING gin (line_name gin_trgm_ops);
CREATE INDEX idx_customers_on_level_and_blacklisted ON customers (level, is_blacklisted);
CREATE INDEX idx_customers_on_last_visit_at ON customers (last_visit_at DESC);

CREATE INDEX idx_staff_users_on_username_trigram ON staff_users USING gin (username gin_trgm_ops);
CREATE INDEX idx_staff_users_on_role_and_is_active ON staff_users (role, is_active);

CREATE INDEX idx_products_on_name_trigram ON products USING gin (name gin_trgm_ops);
CREATE INDEX idx_products_on_store_id_and_brand_id ON products (store_id, brand_id);
CREATE INDEX idx_products_on_store_id_and_category_id ON products (store_id, category_id);
CREATE INDEX idx_products_on_store_id_and_stock_levels ON products (store_id, current_stock, safety_stock);

CREATE INDEX idx_stock_usages_on_product_id_and_is_in_use ON stock_usages (product_id, is_in_use);

CREATE INDEX idx_expenses_on_store_id_and_category ON expenses (store_id, category);
CREATE INDEX idx_expenses_on_store_id_and_supplier_id ON expenses (store_id, supplier_id);
CREATE INDEX idx_expenses_on_store_id_and_payer_id ON expenses (store_id, payer_id);
CREATE INDEX idx_expenses_on_store_id_and_is_reimbursed ON expenses (store_id, is_reimbursed);
