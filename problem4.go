package main

import (
	"log"
	"math"
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

//exp returns a^{pow} mod n
func exp(a, pow, n int) int {
	square_iterations := int(math.Floor(math.Log2(float64(pow))))
	//We need to square a "square_iterations" times
	//then multiply by a^(pow-(2^square_iterations))

	if square_iterations == 0 {
		return 1
	}

	tmp := a
	for i := 0; i < square_iterations; i++ {
		tmp = (tmp * tmp) % n
	}

	//Now, tmp = a^(2^square_iterations)

	//TODO make this more performant
	//return tmp * exp(a, pow - int(math.Pow(2, float64(square_iterations))), n) % n
	for i := 0; i < (pow - int(math.Pow(2, float64(square_iterations)))); i++ {
		tmp = (tmp * a) % n
	}
	return tmp
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

	log.Print(exp(2, 10, 5))
	log.Print(exp(3, 16, 3000))

}
