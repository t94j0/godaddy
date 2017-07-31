package godaddy

type Client struct {
	Key     string
	Secret  string
	Contact Contact
}

func GetName() string {
	return "GoDaddy"
}

func NewClient(key, secret string, contact Contact) *Client {
	if key == "" || secret == "" {
		return nil
	}
	return &Client{key, secret, contact}
}
