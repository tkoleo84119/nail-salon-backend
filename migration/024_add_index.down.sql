DROP INDEX IF EXISTS idx_bookings_on_customer_id;
DROP INDEX IF EXISTS idx_bookings_on_time_slot_id;
DROP INDEX IF EXISTS idx_bookings_on_stylist_id;
DROP INDEX IF EXISTS idx_bookings_on_store_id_and_status;
DROP INDEX IF EXISTS idx_bookings_comprehensive;
DROP INDEX IF EXISTS idx_booking_details_on_booking_id;

DROP INDEX IF EXISTS idx_schedules_on_work_date;
DROP INDEX IF EXISTS idx_schedules_on_store_id_and_stylist_id;

DROP INDEX IF EXISTS idx_time_slots_on_schedule_id;
DROP INDEX IF EXISTS idx_time_slots_on_schedule_id_and_is_available;

DROP INDEX IF EXISTS idx_account_transactions_on_account_id;
DROP INDEX IF EXISTS idx_checkouts_on_booking_id;

DROP INDEX IF EXISTS idx_staff_user_store_access_on_store_id;
DROP INDEX IF EXISTS idx_staff_user_store_access_on_staff_user_id;

DROP INDEX IF EXISTS idx_customers_on_name_trigram;
DROP INDEX IF EXISTS idx_customers_on_phone_trigram;
DROP INDEX IF EXISTS idx_customers_on_line_name_trigram;
DROP INDEX IF EXISTS idx_customers_on_level_and_blacklisted;
DROP INDEX IF EXISTS idx_customers_on_last_visit_at;

DROP INDEX IF EXISTS idx_staff_users_on_username_trigram;
DROP INDEX IF EXISTS idx_staff_users_on_role_and_is_active;

DROP INDEX IF EXISTS idx_products_on_name_trigram;
DROP INDEX IF EXISTS idx_products_on_store_id_and_brand_id;
DROP INDEX IF EXISTS idx_products_on_store_id_and_category_id;
DROP INDEX IF EXISTS idx_products_on_store_id_and_stock_levels;

DROP INDEX IF EXISTS idx_stock_usages_on_product_id_and_is_in_use;

DROP INDEX IF EXISTS idx_expenses_on_store_id_and_category;
DROP INDEX IF EXISTS idx_expenses_on_store_id_and_supplier_id;
DROP INDEX IF EXISTS idx_expenses_on_store_id_and_payer_id;
DROP INDEX IF EXISTS idx_expenses_on_store_id_and_is_reimbursed;

DROP EXTENSION IF EXISTS pg_trgm;
