package json

import (
	"encoding/json"
	"os"
	"testing"
)

func TestEncoder(t *testing.T) {
	writer, _ := os.Create("CustomerOut.json")

	encoder := json.NewEncoder(writer)

	customer := Customer{
		FirstName:  "Di",
		MiddleName: "Han",
		LastName:   "To",
	}
	encoder.Encode(customer)
}
