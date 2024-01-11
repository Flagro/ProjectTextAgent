package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

// Table data structure
type Data struct {
	gorm.Model
	FilePath string
	Text     string
	Metadata string
}

// NewClient creates a new client connected to the PostgreSQL database
func NewClient(host, port, dbName, user, password string) (*Client, error) {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
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
	c.db.Model(&Data{}).Count(&count)
	return count == 0
}

// RemoveData removes data from the database based on the file path
func (c *Client) RemoveData(filePath string) error {
	return c.db.Where("file_path = ?", filePath).Delete(&Data{}).Error
}

// AddData adds data to the database
func (c *Client) AddData(filePath, text, metadata string) error {
	data := Data{
		FilePath: filePath,
		Text:     text,
		Metadata: metadata,
	}
	return c.db.Create(&data).Error
}
