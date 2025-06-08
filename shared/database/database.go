package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cushydigit/nanobank/shared/models"
	_ "github.com/lib/pq"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// init users table
	if err := ensureTable(db, "users", models.CREATE_USERS_TABLE); err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}
	// init accounts table
	if err := ensureTable(db, "accounts", models.CREATE_ACCOUNTS_TABLE); err != nil {
		log.Fatalf("failed to create accounts table: %v", err)
	}
	// init transaction status type
	if err := ensureEnum(db, "transaction_status", models.CREATE_TRANSACTION_STATUS_ENUMS); err != nil {
		log.Fatalf("failed to create enum: %v", err)
	}
	// init transactions table
	if err := ensureTable(db, "transactions", models.CREATE_TRANSACTIONS_TABLE); err != nil {
		log.Fatalf("failed to create transactions table: %v", err)
	}

	return db, nil
}

func ConnectDB(dsn string) *sql.DB {
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

func ensureTable(db *sql.DB, tableName string, createSQL string) error {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM   pg_tables
			WHERE  schemaname = 'public'
			AND    tablename = $1
		);`

	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table %s exists: %w", tableName, err)
	}

	if !exists {
		log.Printf("Table %s does not exist. Creating...", tableName)
		if _, err := db.Exec(createSQL); err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}
		log.Printf("Table %s created successfully.", tableName)
	} else {
		log.Printf("Table %s already exists. Skipping creation.", tableName)
	}

	return nil
}

func ensureEnum(db *sql.DB, enumName string, createSQL string) error {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pg_type
			WHERE typname = $1
		);
	`
	err := db.QueryRow(query, enumName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if enum %s exists: %w", enumName, err)
	}

	if !exists {
		log.Printf("Enum %s does not exist. Creating...", enumName)
		if _, err := db.Exec(createSQL); err != nil {
			return fmt.Errorf("failed to create enum %s: %w", enumName, err)
		}
		log.Printf("Enum %s created successfully.", enumName)
	} else {
		log.Printf("Enum %s already exists. Skipping creation.", enumName)
	}

	return nil
}
