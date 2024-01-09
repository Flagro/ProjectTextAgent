package postgres

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	db *sql.DB
}

func NewClient(url, user, password string) (*Client, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s", user, password, url)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Close() {
	c.db.Close()
}

func (c *Client) IsEmpty() (bool, error) {
	var tableCount int
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"
	err := c.db.QueryRow(query).Scan(&tableCount)
	if err != nil {
		return false, err
	}
	return tableCount == 0, nil
}

func (c *Client) RemoveData(filePath string) error {
	query := "DELETE FROM your_table_name WHERE file_path = $1;"
	_, err := c.db.Exec(query, filePath)
	return err
}

func (c *Client) AddData(filePath, text, metadata string) error {
	// Generate a unique table name (this is a very basic example)
	tableName := fmt.Sprintf("table_%d", rand.Int()) // TODO: fix naming strategy

	// Create table query
	createTableQuery := fmt.Sprintf("CREATE TABLE %s (...);", tableName) // TODO: come up with schema strategy
	_, err := c.db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	// TODO: Insert the CSV data into the table. You will need to parse 'text' and construct an appropriate INSERT query.

	// Associate metadata and filePath with the table
	// Assuming you have a separate metadata table to associate metadata with your data tables
	associateQuery := "INSERT INTO metadata_table (table_name, file_path, metadata) VALUES ($1, $2, $3);"
	_, err = c.db.Exec(associateQuery, tableName, filePath, metadata)
	return err
}
