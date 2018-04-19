package httpclient

// httpclient is a convenience package for executing HTTP
// requests. It's safe in that it always closes response
// bodies and returns byte slices, strings or decodes
// responses into interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Client is an http.Client wrapper
type Client struct {
	Client  *http.Client
	Headers map[string]string
}

// New creates a new Client reference given a client
// timeout
func New(timeout time.Duration, headers map[string]string) *Client {
	return &Client{
		Client:  &http.Client{Timeout: timeout},
		Headers: headers,
	}
}

// Head performs a HEAD request using the passed URL
func (c *Client) Head(url string) error {
	// Execute the request and return the response
	_, err := c.bytes(http.MethodHead, url, nil)
	return err
}

// GetReader performs a GET request using the passed URL
// and returns the io.ReadCloser body
func (c *Client) GetReader(url string) (io.ReadCloser, error) {
	return c.readCloser(http.MethodGet, url, nil)
}

// GetBytes performs a GET request using the passed URL
func (c *Client) GetBytes(url string) ([]byte, error) {
	// Execute the request and return the response
	return c.bytes(http.MethodGet, url, nil)
}

// GetString performs a GET request and returns the response
// as a string
func (c *Client) GetString(url string) (string, error) {
	// Retrieve the bytes and decode the response
	body, err := c.GetBytes(url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetJSON performs a basic http GET request and decodes the JSON
// response into the out interface
func (c *Client) GetJSON(url string, out interface{}) error {
	// Retrieve the bytes and decode the response
	body, err := c.GetBytes(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

// Delete performs a DELETE request using the passed URL
func (c *Client) Delete(url string) error {
	// Execute the request and return the response
	_, err := c.bytes(http.MethodDelete, url, nil)
	return err
}

// bytes executes the passed request using the Client
// http.Client, returning all the bytes read from the response
func (c *Client) bytes(method, url string, in interface{}) ([]byte, error) {
	r, err := c.readCloser(method, url, in)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// readCloser executes the passed request using the Client
// http.Client, returning the io.ReadCloser from the response
func (c *Client) readCloser(method, url string, in interface{}) (io.ReadCloser, error) {
	// Marshal a request body if one exists
	var body io.ReadWriter
	if in != nil {
		if err := json.NewEncoder(body).Encode(in); err != nil {
			return nil, err
		}
	}

	// Generate the request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set all headers
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// Execute the passed request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the status code for an OK
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Non 200 status code : %s", res.Status)
	}

	// Decode and return the bytes
	return res.Body, nil
}
