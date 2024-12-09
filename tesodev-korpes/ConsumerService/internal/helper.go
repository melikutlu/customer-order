package internal

func CalculateVat(price int64) int64 {

	var vatRate int64 = 20
	vatAmount := (price * vatRate) / 100
	totalPrice := price + vatAmount
	return totalPrice
}
