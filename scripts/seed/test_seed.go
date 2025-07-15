package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/snowflake"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		log.Fatal("Please set DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME environment variables")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db := sqlx.MustConnect("pgx", dbURL)
	defer db.Close()

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to initialize snowflake: %v", err)
	}

	// Clear data
	tables := []string{
		"account_transactions", "accounts",
		"stock_usages", "products", "product_categories", "brands",
		"expense_items", "expenses", "suppliers",
		"customer_coupons", "coupons",
		"checkouts", "booking_details", "bookings", "services",
		"time_slot_template_items", "time_slot_templates", "time_slots", "schedules", "stylists",
		"customer_tokens", "customer_auths", "customers",
		"staff_user_tokens", "staff_user_store_access", "staff_users",
		"stores",
	}

	log.Println("🔁 Clearing data...")
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			log.Fatalf("Failed to clear %s: %v", table, err)
		}
	}

	// Create test many stores
	for i := 0; i < 5; i++ {
		storeID := node.Generate().Int64()
		storeName := fmt.Sprintf("測試門市%d", i)
		_, err = db.Exec(`
		INSERT INTO stores (id, name, address, phone)
		VALUES ($1, $2, '台北市信義區', '0223456789')
	`, storeID, storeName)
		if err != nil {
			log.Fatalf("Failed to create test store: %v", err)
		}
	}

	// Create test super admin
	adminID := node.Generate().Int64()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin"), bcrypt.DefaultCost)

	_, err = db.Exec(`
		INSERT INTO staff_users (id, username, email, password_hash, role, is_active)
		VALUES ($1, 'superadmin', 'superadmin@example.com', $2, 'SUPER_ADMIN', true)
	`, adminID, string(hashedPassword))
	if err != nil {
		log.Fatalf("Failed to create super admin: %v", err)
	}

	log.Println("✅ Test seed data created successfully")
}
