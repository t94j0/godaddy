package godaddy

import "errors"

type Client struct {
	Key     string
	Secret  string
	Contact Contact
}

var ErrNoKeySecret = errors.New("Key or secret not supplied")

func (c *Client) GetName() string {
	return "GoDaddy"
}

func NewClient(key, secret string, contact Contact) (*Client, error) {
	if key == "" || secret == "" {
		return nil, ErrNoKeySecret
	}
	return &Client{key, secret, contact}, nil
}
