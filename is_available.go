package godaddy

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// AvailableResponse is the structure that contains how responses are given
// by the GoDaddy API
type AvailableResponse struct {
	Available  bool   `json:"available"`
	Domain     string `json:"domain"`
	Definitive bool   `json:"definitive"`
	Price      uint64 `json:"price"`
	Currency   string `json:"currency"`
	Period     int    `json:"period"`
}

// AvailableRoot is the address for the domain availability API endpoint
const AvailableRoot = "https://api.ote-godaddy.com/v1/domains/available"

// IsAvailable
func (c *Client) IsAvailable(domain string) (bool, uint64, error) {
	// Create GET query parameters
	queryParam := url.Values{}
	queryParam.Set("domain", domain)
	queryParam.Set("checkType", "FAST")
	queryParam.Set("forTransfer", "false")

	// Generate the URL given the query parameters
	reqURL, _ := url.Parse(AvailableRoot)
	reqURL.RawQuery = queryParam.Encode()

	// Generate the Authorization header with the Client
	authHeader := "sso-key " + c.Key + ":" + c.Secret

	// Create the HTTP client used for making the request
	httpClient := http.DefaultClient

	// Create response object with given data
	request, err := http.NewRequest("GET", reqURL.String(), nil)
	request.Header.Add("Authorization", authHeader)
	if err != nil {
		return false, 0, err
	}

	url, err := http.ProxyFromEnvironment(request)
	if err != nil {
		return false, 0, err
	}

	response.URL = url

	// make request
	response, err := httpClient.Do(request)
	if err != nil {
		return false, 0, err
	}

	// Parse request
	var responseObj AvailableResponse

	if err := json.NewDecoder(response.Body).Decode(&responseObj); err != nil {
		return false, 0, err
	}

	return responseObj.Available, responseObj.Price, nil
}
