package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONArray(t *testing.T) {
	customer := Customer{
		FirstName:  "Di",
		MiddleName: "Han",
		LastName:   "To",
		Hobbies:    []string{"Gaming", "Coding", "Hiking"},
	}
	bytes, _ := json.Marshal(customer)
	fmt.Println(string(bytes))
}

func TestJSONArrayDecode(t *testing.T) {
	jsonString := `{"FirstName":"Di","MiddleName":"Han","LastName":"To","Age":0,"Married":false,"Hobbies":["Gaming","Coding","Hiking"],"Addresses":null}`
	jsonBytes := []byte(jsonString)

	customer := &Customer{}
	json.Unmarshal(jsonBytes, customer)

	fmt.Println(customer)
	fmt.Println(customer.FirstName)
	fmt.Println(customer.Hobbies)
}

func TestJSONArrayComplex(t *testing.T) {
	customer := Customer{
		FirstName: "Di",
		Addresses: []Address{
			{
				Street:     "Jalan Jalan Jalan",
				Country:    "Indonesia",
				PostalCode: "46471",
			},
			{
				Street:     "Jalan Masih Jelek",
				Country:    "Indonesia",
				PostalCode: "46472",
			},
		},
	}
	bytes, _ := json.Marshal(customer)

	fmt.Println(string(bytes))
}

func TestJSONArrayComplexDecode(t *testing.T) {
	jsonString := `{"FirstName":"Di","MiddleName":"","LastName":"","Age":0,"Married":false,"Hobbies":null,"Addresses":[{"Street":"Jalan Jalan Jalan","Country":"Indonesia","PostalCode":"46471"},{"Street":"Jalan Masih Jelek","Country":"Indonesia","PostalCode":"46472"}]}`
	jsonBytes := []byte(jsonString)

	customer := &Customer{}

	json.Unmarshal(jsonBytes, customer)

	fmt.Println(customer)
	fmt.Println(customer.FirstName)
	fmt.Println(customer.Addresses)
}

func TestOnlyJSONArray(t *testing.T) {
	addresses := []Address{
		{
			Street:     "Jalan Jalan Jalan",
			Country:    "Indonesia",
			PostalCode: "46471",
		},
		{
			Street:     "Jalan Masih Jelek",
			Country:    "Indonesia",
			PostalCode: "46472",
		},
	}
	bytes, _ := json.Marshal(addresses)

	fmt.Println(string(bytes))
}
