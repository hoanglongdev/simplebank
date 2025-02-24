package util

const (
	USD = "USD"
	VND = "VND"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, VND:
		return true
	}
	return false
}
