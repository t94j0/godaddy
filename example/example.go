package main

import (
	"fmt"
	"os"

	"github.com/t94j0/godaddy"
)

const Key = ""
const Secret = ""

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [domain]", os.Args[0])
		return
	}

	dmn := os.Args[1]

	client := godaddy.NewClient(
		Key,
		Secret,
		godaddy.Contact{
			"Max",
			"",
			"Harley",
			"Max Co.",
			"CEO",
			"maxh@maxh.io",
			"+1.9999999999",
			"",
			godaddy.Address{
				"12 Awesome Blvd",
				"Charleston",
				"SC",
				"29463",
				"US",
			},
		},
	)

	fmt.Println("Purchasing domain:", dmn)

	isAvail, _, err := client.IsAvailable(dmn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isAvail {
		if err := client.Purchase(dmn); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Domain is unavailable.")
	}
}
