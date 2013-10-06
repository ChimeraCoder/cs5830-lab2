package main

import (
	"math/big"
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
func Exp(a, pow, n big.Int) big.Int {

	result := big.NewInt(1)
	tmp := &a

	for i := 0; i < pow.BitLen(); i++ {
		bit := pow.Bit(i)
		if bit == 1 {
			result.Mul(result, tmp)
			result.Mod(result, &n)
		}

		//This could be simplified by big.Exp, but we can't use that
		tmp.Mul(tmp, tmp)
		tmp.Mod(tmp, &n)
	}
	return *result
}

//invert finds the modular inverse of an element, mod divisor
func invert(element, divisor int) int {
	g := gcd(element, divisor)
	s, _ := euclid(element, divisor)

	//The pair (s/g, t/g) is the solution to ax + my = 1
	//where t is the discarded return value from euclid()
	return s / g
}

//MillerRabin checks if the number is prime, using the Miller-Rabin test
//It may return false positives (ie, incorrectly identify a composite number as prime)
//It will never incorrectly reject a prime number as composite
//Higer values of numTests will decrease the chance of a false positive
func MillerRabin(n big.Int, numTests int) bool {
	for i := 0; i < numTests; i++ {
		if MillerRabinAux(n) == false {
			return false
		}
	}
	return true
}

func MillerRabinAux(n big.Int) bool {
	d := big.NewInt(0)
	d = d.Sub(&n, big.NewInt(1))

	s := big.NewInt(0)

	result := big.NewInt(0)

	for {
		if result := result.Mod(d, big.NewInt(2)); result.Cmp(big.NewInt(0)) == 0 {
			break
		}

		d = d.Div(d, big.NewInt(2)) //d = d/2
		s.Add(s, big.NewInt(1))     // s++
	}

	//Without math.Big, the following would look like
	//a := r.Int63n(n-3) + 2 // returns vals between [2, n-2] incl.
	a := big.NewInt(0)
	upper := big.NewInt(0)
	upper = upper.Sub(&n, big.NewInt(3))
	a = a.Rand(r, upper)
	a = a.Add(a, big.NewInt(2))

	x := Exp(*a, *d, n)

	//Equivalent to
	//if (x == 1) || (x == n-1) {
	f := big.NewInt(0)
	if (x.Cmp(big.NewInt(1)) == 0) || (x.Cmp(f.Sub(&n, big.NewInt(1))) == 0) {
		return true
	}

	s_minus_one := big.NewInt(0)
	s_minus_one = s_minus_one.Sub(s, big.NewInt(1))

	for i := big.NewInt(0); i.Cmp(s_minus_one) == -1; i.Add(i, big.NewInt(1)) {
		x = Exp(x, *big.NewInt(2), n)
		if x.Cmp(big.NewInt(1)) == 0 {
			return false
		}

		tmp := big.NewInt(0)
		if x.Cmp(tmp.Sub(&n, big.NewInt(1))) == 0 {
			return true
		}
	}
	return false
}

//RandomNBitNumber returns a random number with the specified number of bits
func RandomNBitNumber(n int64) big.Int {

	bigN := big.NewInt(n - 1)
	bigN = bigN.Exp(big.NewInt(2), big.NewInt(n), nil)

	result := big.NewInt(0)
	result = result.Rand(r, bigN)

	result.Add(result, bigN)
	return *result
}

//RandomNBitPrime returns random prime numbers of the specified size
func RandomNBitPrime(n int64, certainty int) big.Int {
	for {
		result := RandomNBitNumber(n)
		if MillerRabin(result, certainty) {
			return result
		}
	}
}

//RandomNBitSafePrime is like RandomNBitPrime, except it is guaranteed to return
//only safe primes
func RandomNBitSafePrime(n int64, certainty int) big.Int {
	for {
		number := RandomNBitPrime(n, certainty)
		other := big.NewInt(0)
		other = other.Sub(&number, big.NewInt(1))
		other = other.Div(other, big.NewInt(2))
		if MillerRabin(*other, certainty) {
			return number
		}
	}
}

func main() {
}
