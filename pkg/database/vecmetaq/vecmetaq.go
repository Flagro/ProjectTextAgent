package vecmetaq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// VecMetaQClient holds the configuration for the VecMetaQ database API.
type Client struct {
	Host       string
	Port       string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewClient creates a new VecMetaQClient with the necessary configuration.
func NewClient(host, port, username, password string) (*Client, error) {
	return &Client{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{},
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

// PostText posts text, tag, and metadata to the VecMetaQ database.
func (c *Client) PostText(text, tag string, metadata map[string]interface{}) error {
	endpoint := fmt.Sprintf("%s:%s/add_data/", c.Host, c.Port)
	requestBody, err := json.Marshal(map[string]interface{}{
		"text":     text,
		"tag":      tag,
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Handle non-OK response
		return err
	}

	return nil
}

// DeleteTag deletes a tag from the VecMetaQ database.
func (c *Client) DeleteTag(tag string) error {
	endpoint := fmt.Sprintf("%s:%s/delete_data/", c.Host, c.Port)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	query := req.URL.Query()
	query.Add("tag", tag)
	req.URL.RawQuery = query.Encode()
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Handle non-OK response
		return err
	}

	return nil
}
