package services

import (
	"strconv"
)

func parsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}
