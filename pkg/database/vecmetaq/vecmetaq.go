package vecmetaq

import (
	"bytes"
	"encoding/json"
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

// PostText posts text, tag, and metadata to the VecMetaQ database.
func (c *Client) PostText(text, tag string, metadata map[string]interface{}) error {
	endpoint := c.BaseURL + "/path/to/post/endpoint"
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
	endpoint := c.BaseURL + "/path/to/delete/endpoint"
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	req.URL.Query().Add("tag", tag)
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
