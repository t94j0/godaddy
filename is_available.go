package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type AvailableResponse struct {
	Available  bool   `json:"available"`
	Domain     string `json:"domain"`
	Definitive bool   `json:"definitive"`
	Price      uint64 `json:"price"`
	Currency   string `json:"currency"`
	Period     int    `json:"period"`
}

const GoDaddyRoot = "https://api.ote-godaddy.com/v1/domains/available"

func IsAvailable(domain string) (bool, uint64, error) {
	queryParam := url.Values{}
	queryParam.Set("domain", domain)
	queryParam.Set("checkType", "FAST")
	queryParam.Set("forTransfer", "false")

	url := url.Parse(GoDaddyRoot)
	url.Query = queryParam

	response, err := http.Get(url.String())
	if err != nil {
		return false, 0, err
	}

	var responseObj AvailableResponse

	if err := json.NewDecoder(response.Body).Decode(&responseObj); err != nil {
		return false, 0, err
	}

	return responseObj.Available, responseObj.Price, nil
}
