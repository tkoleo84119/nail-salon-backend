seed-test:
	@echo "Seeding test data..."
	go run scripts/seed/test_seed.go

migrate-up:
	@echo "Migrating up..."
	migrate -path ./migration -database "$(DB_URL)" up

migrate-down:
	@echo "Migrating down..."
	migrate -path ./migration -database "$(DB_URL)" down $(NUMBER)