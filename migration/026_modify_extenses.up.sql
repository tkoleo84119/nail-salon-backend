ALTER TABLE expenses
ADD COLUMN updater bigint;

ALTER TABLE expenses
ADD CONSTRAINT fk_expenses_updater
FOREIGN KEY (updater)
REFERENCES staff_users(id)
ON DELETE CASCADE;