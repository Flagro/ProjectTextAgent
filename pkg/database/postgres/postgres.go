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
	FilePath   string
	Metadata   string
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
func (c *Client) IsEmpty() (bool, error) {
	var count int64
	if err := c.db.Model(&TableData{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// RemoveData removes the data from the database
func (c *Client) RemoveData(filePath string) error {
	if err := c.db.Where("file_path = ?", filePath).Delete(&TableData{}).Error; err != nil {
		return err
	}
	return nil
}

// AddData adds the data to the database
func (c *Client) AddData(filePath, text, metadata string) error {
	entry := TableData{FilePath: filePath, CSVContent: text, Metadata: metadata}
	if err := c.db.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}
