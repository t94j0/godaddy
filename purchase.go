package godaddy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/publicsuffix"
)

// LegalAgreement is an object for gettin the Consent key
type LegalAgreementResponse struct {
	AgreementKey string `json:"agreementKey,omitempty"`
	Title        string `json:"title,omitempty"`
	Url          string `json:"url,omitempty"`
	Content      string `json:"content,omitempty"`
	// Used for errors
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Address is used for the Contact object
type Address struct {
	Address1   string `json:"address1,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Contact is what is listed in the whois information if privact isn't selected
type Contact struct {
	First          string  `json:"nameFirst,omitempty"`
	Middle         string  `json:"nameMiddle,omitempty"`
	Last           string  `json:"nameLast,omitempty"`
	Organization   string  `json:"organization,omitempty"`
	JobTitle       string  `json:"jobTitle,omitempty"`
	Email          string  `json:"email,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	Fax            string  `json:"fax,omitempty"`
	AddressMailing Address `json:"addressMailing,omitempty"`
}

// Consent is needed for TLDs that need to be purchased
type Consent struct {
	AgreementKeys []string `json:"agreementKeys,omitempty"`
	AgreedBy      string   `json:"agreedBy,omitempty"`
	AgreedAt      string   `json:"agreedAt,omitempty"`
}

// DomainPurchase is the object needed for purchasing a domain
type DomainPurchase struct {
	Domain            string   `json:"domain,omitempty"`
	Consent           Consent  `json:"consent,omitempty"`
	Period            int32    `json:"period,omitempty"`
	NameServers       []string `json:"nameServers,omitempty"`
	RenewAuto         bool     `json:"renewAuto,omitempty"`
	Privacy           bool     `json:"privacy,omitempty"`
	ContactRegistrant Contact  `json:"contactRegistrant,omitempty"`
	ContactAdmin      Contact  `json:"contactAdmin,omitempty"`
	ContactTech       Contact  `json:"contactTech,omitempty"`
	ContactBilling    Contact  `json:"contactBilling,omitempty"`
}

// DomainPurchaseResponse is the response object from buying a domain
type DomainPurchaseResponse struct {
	// The actual response
	OrderID   int32  `json:"orderId,omitempty"`
	ItemCount int32  `json:"itemCount,omitempty"`
	Total     int32  `json:"total,omitempty"`
	Currency  string `json:"currency,omitempty"`

	// Error
	Code    string `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

// AgreementRoot is the location for agreeing to purchasing under TLD
const AgreementRoot = "https://api.ote-godaddy.com/v1/domains/agreements"

// PurchaseRoot location for purchasing domain
const PurchaseRoot = "https://api.ote-godaddy.com/v1/domains/purchase"

// ErrPurchasing is returned if there is an error purchasing the domain
var ErrPurchasing = errors.New("Error purchasing the domain")

// Purchase handles the purchasing of the specified domain. For GoDaddy, the
// process for purchasing a domain is:
// 1. Get consent key from the AgreementRoot by creating a LegalAgreement
//    object
// 2. Create DomainPurchase object
// 3. Purchase domain by making a request to PurchaseRoot
func (c *Client) Purchase(domain string) error {
	privacy := true

	// First, we need to create a Consent object
	// Parse domain for TLD
	tld, _ := publicsuffix.PublicSuffix(domain)

	// Create query parameters
	query := url.Values{}
	query.Set("tlds", tld)
	if privacy {
		query.Set("privacy", "true")
	} else {
		query.Set("privacy", "false")
	}

	// Make HTTP client
	client := http.DefaultClient

	// Generate auth header
	authHeader := "sso-key " + c.Key + ":" + c.Secret

	// Make request to the agreement page
	request, err := http.NewRequest("GET", AgreementRoot+"?"+query.Encode(), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", authHeader)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// Get the agreement key for the TLD
	var legal []LegalAgreementResponse
	if err := json.NewDecoder(response.Body).Decode(&legal); err != nil {
		return err
	}

	var consent Consent
	keys := make([]string, 0)
	for _, agreement := range legal {
		keys = append(keys, agreement.AgreementKey)
	}
	consent.AgreementKeys = keys
	consent.AgreedBy = c.Contact.First + " " + c.Contact.Last
	consent.AgreedAt = time.Now().Format("2006-01-02T15:04:05-0700Z")

	// Next, create object for purchasing the domain
	// TODO: Allow users to create this object
	var purchase DomainPurchase
	purchase.Domain = domain
	purchase.Consent = consent
	purchase.Period = 1
	purchase.NameServers = []string{}
	purchase.RenewAuto = false
	purchase.Privacy = privacy
	purchase.ContactRegistrant = c.Contact
	purchase.ContactAdmin = c.Contact
	purchase.ContactTech = c.Contact
	purchase.ContactBilling = c.Contact

	// Create buffer for the body
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(purchase); err != nil {
		return err
	}

	fmt.Println(body.String())

	// Lastly, purchase the domain
	fmt.Println(PurchaseRoot)
	request, err = http.NewRequest("POST", PurchaseRoot, body)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", authHeader)
	request.Header.Add("Content-Type", "application/json")

	fmt.Println("")

	response, err = client.Do(request)
	if err != nil {
		return err
	}

	// Marshal purchase response
	var dpResp DomainPurchaseResponse

	if err := json.NewDecoder(response.Body).Decode(&dpResp); err != nil {
		return err
	}

	if dpResp.Code != "" {
		fmt.Fprintf(os.Stderr, "%+v\n", dpResp)
		return ErrPurchasing
	}

	return nil
}
