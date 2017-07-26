package godaddy

type Client struct {
	Key    string
	Secret string
}

func NewClient(key, secret string) *Client {
	return &Client{key, secret}
}
