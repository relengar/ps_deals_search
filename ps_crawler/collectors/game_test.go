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

func TestGetPrice(t *testing.T) {
	priceTexts := []struct {
		text     string
		expected float32
		isNum    bool
	}{
		{
			text:     "$123.456",
			expected: 123.456,
			isNum:    true,
		},
		{
			text:     "not valid",
			expected: 0,
			isNum:    false,
		},
	}

	for _, v := range priceTexts {
		price, err := getPrice(v.text)
		if !v.isNum {
			assert.NotNil(t, err)
			return
		}

		assert.Nil(t, err)
		assert.Equal(t, v.expected, price)
	}
}

func TestGetCurrency(t *testing.T) {
	priceText := "$123.456"
	currency := getCurrency(priceText)
	assert.Equal(t, "$", currency)
}
