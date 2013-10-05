package main

import (
	"log"
	"math"
    "math/rand"
)

// using constant for debugging; in production would use time.now()
var r = rand.New(rand.NewSource(123))

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
func Exp(a, pow, n int) int {
	if pow  == 0 {
		return 1
	}

	//We need to square a "square_iterations" times
	//then multiply by a^(pow-(2^square_iterations))
	square_iterations := int(math.Floor(math.Log2(float64(pow))))

	tmp := a
	for i := 0; i < square_iterations; i++ {
		tmp = (tmp * tmp) % n
	}

	//Now, tmp = a^(2^square_iterations)

	//TODO make this more performant
	return tmp * Exp(a, pow - int(math.Pow(2, float64(square_iterations))), n) % n
}

//invert finds the modular inverse of an element, mod divisor
func invert(element, divisor int) int {
	g := gcd(element, divisor)
	s, _ := euclid(element, divisor)

	//The pair (s/g, t/g) is the solution to ax + my = 1
	//where t is the discarded return value from euclid()
	return s / g
}

func MillerRabin(n, numTests int) bool {
    for i := 0; i < numTests; i++ {
        if (MillerRabinAux(n) == false) {
            return false
        }
    }
    return true
}

func MillerRabinAux(n int) bool {
    d := n-1
    s := 1
    for ; d%s == 0; {
        s *= 2
    }
    s = int(math.Floor(math.Log2(float64(s))))
    d = (n-1)/s
    // crufty
    //s := int(math.Floor(math.Log2(float64(n-1))))
    //d := (n-1) / int(math.Pow(2,float64(s)))
    a := r.Intn(n-3) + 2 // returns vals between [2, n-2] incl.
    x := Exp(a, d, n)
    if (x == 1) || (x == n-1) {
        return true
    }
    for i := 0; i < s-1; i++ {
        x = Exp(x, 2, n)
        if (x == 1) {
            return false
        }
        if (x == n-1) {
            continue
        }
    }
    return false
}

func main() {
	log.Print(gcd(10, 15))
	log.Print(euclid(10, 15))
	log.Print(invert(101, 102))

	log.Print(Exp(2, 10, 5))
	log.Print(Exp(3, 16, 3000))

}
