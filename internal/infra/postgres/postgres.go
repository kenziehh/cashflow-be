package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kenziehh/cashflow-be/config"
)

func InitDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Run migrations
	RunMigrations(db)

	log.Println("Database connected successfully")
	return db
}


func RunMigrations(db *sql.DB) {
	migrationsDir := "database/migrations" // relative to /app in container

	// Ensure schema_migrations table exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		log.Fatalf("failed to ensure schema_migrations table: %v", err)
	}

	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("failed to read migrations dir: %v", err)
	}

	// Collect .sql files with version prefix
	type migrationFile struct {
		Version int
		Name    string
		Path    string
	}

	var migrations []migrationFile

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(f.Name(), "_", 2)
		if len(parts) < 2 {
			log.Printf("skip file without version prefix: %s", f.Name())
			continue
		}

		v, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("skip file with invalid version prefix: %s", f.Name())
			continue
		}

		migrations = append(migrations, migrationFile{
			Version: v,
			Name:    f.Name(),
			Path:    filepath.Join(migrationsDir, f.Name()),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Run each migration if not applied
	for _, m := range migrations {
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, m.Version).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to check migration %d: %v", m.Version, err)
		}

		if exists {
			log.Printf("migration %03d already applied, skipping (%s)", m.Version, m.Name)
			continue
		}

		content, err := ioutil.ReadFile(m.Path)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", m.Path, err)
		}

		sqlText := string(content)
		log.Printf("applying migration %03d: %s", m.Version, m.Name)

		if _, err := db.Exec(sqlText); err != nil {
			log.Fatalf("failed to execute migration %s: %v", m.Name, err)
		}

		if _, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, m.Version); err != nil {
			log.Fatalf("failed to record migration %d: %v", m.Version, err)
		}
	}

	log.Println("✅ Custom migrations completed")
}
