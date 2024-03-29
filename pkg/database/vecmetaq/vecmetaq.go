package vecmetaq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// VecMetaQClient holds the configuration for the VecMetaQ database API.
type Client struct {
	Endpoint      string
	Username      string
	Password      string
	HTTPClient    *http.Client
	Retries       int
	RetryInterval time.Duration
}

// NewClient creates a new VecMetaQClient with the necessary configuration.
func NewClient(host, port, username, password string) (*Client, error) {
	return &Client{
		Endpoint:      fmt.Sprintf("http://%s:%s", host, port),
		Username:      username,
		Password:      password,
		HTTPClient:    &http.Client{},
		Retries:       3,
		RetryInterval: time.Second * 5,
	}, nil
}

func (c *Client) Close() error {
	return nil
}

// IsEmpty checks if the database is empty
// TODO:
// currently there are no interfaces to check if this is false, so we will have to return true
func (c *Client) IsEmpty() (bool, error) {
	return true, nil
}

// AddData posts text, filePath, and metadata to the VecMetaQ database.
func (c *Client) AddData(filePath, text, metadata string) error {
	endpoint := fmt.Sprintf("%s/add_data/", c.Endpoint)
	requestBody, err := json.Marshal(map[string]interface{}{
		"text":     text,
		"tag":      filePath,
		"metadata": metadata,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	return c.retriesDo(req)
}

// RemoveData deletes a tag from the VecMetaQ database.
func (c *Client) RemoveData(filePath string) error {
	endpoint := fmt.Sprintf("%s/delete_data/", c.Endpoint)
	requestBody, err := json.Marshal(map[string]interface{}{
		"tag": filePath,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	return c.retriesDo(req)
}

func (c *Client) retriesDo(req *http.Request) error {
	for i := 0; i < c.Retries; i++ {
		connection_err, status_err := c.do(req)
		if connection_err == nil {
			return status_err
		}
		if i == c.Retries-1 {
			return connection_err
		}
		log.Println("Error establishin connection with VecMetaQ, retrying:", connection_err)
		time.Sleep(c.RetryInterval)
	}
	return nil
}

// Do request with retries wrapper
func (c *Client) do(req *http.Request) (error, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Handle non-OK response
		return nil, err
	}
	return nil, nil
}
