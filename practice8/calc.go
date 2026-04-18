package main
import "fmt"

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("bratan, nelzya delit na nol")
	}
	return a / b, nil
}