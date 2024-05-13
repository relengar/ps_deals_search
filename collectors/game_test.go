package collectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpirationParser(t *testing.T) {
	expirations := []struct {
		text  string
		valid bool
	}{
		{
			text:  "Offer ends 23/5/2024 00:59 CEST",
			valid: true,
		},
		{
			text:  "Offer ends 5/23/2024 01:59 PM UTC",
			valid: true,
		},
		{
			text:  "Offer ends at 5/23/2024 01:59 PM UTC",
			valid: false,
		},
	}

	for _, v := range expirations {
		date, err := toExpiration(v.text)
		if v.valid {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
		fmt.Println(date)
	}
}
