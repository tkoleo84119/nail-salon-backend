ALTER TABLE expenses
DROP CONSTRAINT fk_expenses_updater;

ALTER TABLE expenses
DROP COLUMN updater;
