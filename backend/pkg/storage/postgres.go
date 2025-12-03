package storage

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
)

// PostgresRepository implements Repository using PostgreSQL and JSONB.
type PostgresRepository[T any] struct {
	db         *sql.DB
	tableName  string
	collection string
}

// NewPostgresRepository creates a new PostgresRepository.
// It ensures the table exists.
func NewPostgresRepository[T any](db *sql.DB, collection string) *PostgresRepository[T] {
	log.Printf("DEBUG: Creating PostgresRepository for collection: %s", collection)
	repo := &PostgresRepository[T]{
		db:         db,
		tableName:  "items",
		collection: collection,
	}
	repo.ensureTable()
	return repo
}

func (r *PostgresRepository[T]) ensureTable() {
	query := `
	CREATE TABLE IF NOT EXISTS items (
		id TEXT NOT NULL,
		collection TEXT NOT NULL,
		data JSONB NOT NULL,
		PRIMARY KEY (id, collection)
	);
	CREATE INDEX IF NOT EXISTS idx_items_collection ON items(collection);
	`
	if _, err := r.db.Exec(query); err != nil {
		log.Printf("Error creating table: %v", err)
	}
}

// Save stores the value as JSONB.
func (r *PostgresRepository[T]) Save(id string, value T) {
	log.Printf("DEBUG: Saving item %s to collection %s", id, r.collection)
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	query := `
	INSERT INTO items (id, collection, data)
	VALUES ($1, $2, $3)
	ON CONFLICT (id, collection) DO UPDATE SET
		data = EXCLUDED.data;
	`
	if _, err := r.db.Exec(query, id, r.collection, data); err != nil {
		log.Printf("Error saving item %s: %v", id, err)
	} else {
		log.Printf("DEBUG: Successfully saved item %s", id)
	}
}

// Get retrieves an item by id.
func (r *PostgresRepository[T]) Get(id string) (T, bool) {
	var data []byte
	var zero T

	query := `SELECT data FROM items WHERE id = $1 AND collection = $2`
	err := r.db.QueryRow(query, id, r.collection).Scan(&data)
	if err == sql.ErrNoRows {
		return zero, false
	}
	if err != nil {
		log.Printf("Error getting item %s: %v", id, err)
		return zero, false
	}

	if err := json.Unmarshal(data, &zero); err != nil {
		log.Printf("Error unmarshaling data: %v", err)
		return zero, false
	}

	return zero, true
}

// List returns all values for the collection.
func (r *PostgresRepository[T]) List() []T {
	log.Printf("DEBUG: Listing items for collection %s", r.collection)
	query := `SELECT data FROM items WHERE collection = $1`
	rows, err := r.db.Query(query, r.collection)
	if err != nil {
		log.Printf("Error listing items: %v", err)
		return nil
	}
	defer rows.Close()

	var values []T
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		var v T
		if err := json.Unmarshal(data, &v); err != nil {
			log.Printf("Error unmarshaling row: %v", err)
			continue
		}
		values = append(values, v)
	}
	log.Printf("DEBUG: Found %d items for collection %s", len(values), r.collection)
	return values
}
