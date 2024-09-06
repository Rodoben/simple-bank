package util

const (
	INR = "INR"
	USD = "USD"
	AED = "AED"
	TBH = "TBH"
)

func IsSuppportedCurrency(currency string) bool {
	switch currency {
	case INR, USD, AED, TBH:
		return true
	}
	return false
}
