package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"math/rand"
)

func SeedCategories(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	for _, name := range categories {
		id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

		query := `INSERT INTO categories (id, name)
		          VALUES ($1, $2)
		          ON CONFLICT (id) DO NOTHING;`

		_, err := db.ExecContext(ctx, query, id, name)
		if err != nil {
			log.Printf("❌ Failed to seed category '%s': %v", name, err)
			continue
		}

		fmt.Printf("✅ Seeded category: %s\n", name)
	}

	return nil
}
