package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func SeedCategoriesIfEmpty(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check categories count: %w", err)
	}

	if count > 0 {
		fmt.Println("⚠️  Categories table already has data, skipping seeding...")
		return nil
	}

	categories := []string{
		"Food & Drinks",
		"Transportation",
		"Utilities",
		"Entertainment",
		"Health",
		"Education",
		"Shopping",
		"Savings",
		"Investment",
		"Others",
	}

	fmt.Println("🌱 Seeding categories...")

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, name := range categories {
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

		query := `INSERT INTO categories (id, name)
		          VALUES ($1, $2)
		          ON CONFLICT (name) DO NOTHING;`

		_, err := db.ExecContext(ctx, query, id, name)
		if err != nil {
			log.Printf("❌ Failed to seed category '%s': %v", name, err)
			continue
		}

		fmt.Printf("✅ Seeded category: %s\n", name)
	}

	fmt.Println("🎉 Categories seeding completed successfully!")
	return nil
}
