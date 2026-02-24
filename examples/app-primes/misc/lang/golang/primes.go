package main

import (
	"fmt"
	"os"
	"strconv"
)

func isPrime(n int) bool {
	ret := n > 1
	for i := 2; i < n; i++ {
		if n%i == 0 {
			ret = false
		}
	}
	return ret
}

func getNumPrimes(start int, len int) int {
	ret := 0
	for k := 0; k < len; k++ {
		if isPrime(start + k) {
			ret++
		}
	}
	return ret
}

func main() {
	start, _ := strconv.Atoi(os.Args[1])
	len, _ := strconv.Atoi(os.Args[2])

	fmt.Println(getNumPrimes(start, len))
}
