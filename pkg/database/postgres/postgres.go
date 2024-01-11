package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

// Table data structure
type TableData struct {
	gorm.Model
	FilePath   string `gorm:"primary_key"`
	Metadata   string `gorm:"primary_key"`
	CSVContent string
}

// NewClient creates a new PostgreSQL client
func NewClient(host, port, dbName, user, password string) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbName, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&TableData{})

	return &Client{db: db}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// IsEmpty checks if the database is empty
func (c *Client) IsEmpty() bool {
	var count int64
	c.db.Model(&TableData{}).Count(&count)
	return count == 0
}

// RemoveData removes the data from the database
func (c *Client) RemoveData(filePath string) {
	c.db.Delete(&TableData{}, "file_path = ?", filePath)
}

// AddData adds the data to the database
func (c *Client) AddData(filePath, text, metadata string) {
	c.db.Create(&TableData{FilePath: filePath, CSVContent: text, Metadata: metadata})
}
