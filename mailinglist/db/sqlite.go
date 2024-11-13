package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db := &DB{Conn: conn}
	if err := db.createTable(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS subscribers (email TEXT PRIMARY KEY);`
	_, err := db.Conn.Exec(query)
	return err
}

func (db *DB) AddSubscriber(email string) error {
	_, err := db.Conn.Exec("INSERT INTO subscribers (email) VALUES (?)", email)
	return err
}

func (db *DB) ListSubscribers() ([]string, error) {
	rows, err := db.Conn.Query("SELECT email FROM subscribers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}
	return emails, nil
}
