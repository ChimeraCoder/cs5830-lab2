package main

import (
	"log"
)

//Function gcd uses Euclid's algorithm to compute the inverse of a number, mod m
//ax + my = 1
func gcd(a, m int) int {
	for {
		if m == 0 {
			break
		}
		tmp := m
		m = a % tmp
		a = tmp
	}
	return a
}

func main() {

	//This should be 3
	log.Print(gcd(10, 15))
}
