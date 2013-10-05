package main

import (
	"log"
)

//Function gcd uses Euclid's algorithm to compute the inverse of a number, mod m
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

//euclid returns the numbers x, y such that ax + by = gcd(a, b)
func euclid(a, b int) (x int, y int) {
	for {
		if b == 0 {
			return 1, 0
		}
		q := a / b
		r := a % b
		s, t := euclid(b, r)
		return t, s - q*t
	}
}

//invert finds the modular inverse of an element, mod divisor
func invert(element, divisor int) int {
	g := gcd(element, divisor)
	s, _ := euclid(element, divisor)

	//The pair (s/g, t/g) is the solution to ax + my = 1
	//where t is the discarded return value from euclid()
	return s / g
}

func main() {
	log.Print(gcd(10, 15))
	log.Print(euclid(10, 15))
	log.Print(invert(101, 102))
}
