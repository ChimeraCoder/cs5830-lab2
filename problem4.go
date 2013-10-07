package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"sync"
)

// using constant for debugging; in production would use time.now()
var r = rand.New(rand.NewSource(123))

//Function gcd uses Euclid's algorithm to compute the inverse of a number, mod m
func gcd(a, m *big.Int) *big.Int {
	a2 := big.NewInt(0)
	m2 := big.NewInt(0)
	a2.Set(a)
	m2.Set(m)

	for {
		if m2.Cmp(big.NewInt(0)) == 0 {
			break
		}

		tmp := big.NewInt(0)
		tmp.Set(m2)     //tmp = m
		m2.Mod(a2, tmp) //m = a % tmp
		a2.Set(tmp)     //a = tmp
	}
	return a2
}

//euclid returns the numbers x, y such that ax + by = gcd(a, b)
func euclid(a, b *big.Int) (x *big.Int, y *big.Int) {
	//Copy over a and b into new places in memory
	//So that the caller's values don't get modified
	a2 := big.NewInt(0)
	b2 := big.NewInt(0)
	a2.Set(a)
	b2.Set(b)
	a = a2
	b = b2

	for {
		if b.Cmp(big.NewInt(0)) == 0 {
			return big.NewInt(1), big.NewInt(0)
		}
		q := big.NewInt(0)
		r := big.NewInt(0)

		q.Div(a, b) //q := a / b
		r.Mod(a, b) //r := a % b

		s, t := euclid(b, r)

		//secondVal = s - q*t
		secondVal := big.NewInt(0)
		secondVal.Mul(q, t)
		secondVal.Sub(s, secondVal)

		return t, secondVal
	}
}

//exp returns a^{pow} mod n
func Exp(a, pow, n big.Int) big.Int {

	result := big.NewInt(1)

	//Set tmp to a so that when we take its address, we do not modify a
	tmp := a

	for i := 0; i < pow.BitLen(); i++ {
		bit := pow.Bit(i)
		if bit == 1 {
			result.Mul(result, &tmp)
			result.Mod(result, &n)
		}

		//This could be simplified by big.Exp, but we can't use that
		foo := big.NewInt(0)
		tmp = *foo.Mul(&tmp, &tmp)
		tmp = *foo.Mod(&tmp, &n)
	}

	return *result
}

//invert finds the modular inverse of an element, mod divisor
func invert(element, divisor *big.Int) *big.Int {
	g := gcd(element, divisor)

	s, _ := euclid(element, divisor)

	//The pair (s/g, t/g) is the solution to ax + my = 1
	//where t is the discarded return value from euclid()
	result := big.NewInt(0)

	return result.Div(s, g) //return s/g
}

//MillerRabin checks if the number is prime, using the Miller-Rabin test
//It may return false positives (ie, incorrectly identify a composite number as prime)
//It will never incorrectly reject a prime number as composite
//Higer values of numTests will decrease the chance of a false positive
func MillerRabin(n big.Int, numTests int) bool {

	var wg sync.WaitGroup

	results := make(chan bool)

	for i := 0; i < numTests; i++ {
		n2 := big.NewInt(0)
		n2.Set(&n)
		go func(n2 big.Int) {
			wg.Add(1)
			millerRabinAux(n2, results)
			wg.Done()
		}(*n2)
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	for {
		select {
		case r := <-results:
			{
				if r == false {
					return false
				}
			}

		case <-done:
			{
				return true
			}
		}
	}

	//This should never be reachable
	return false
}

//millerRabinAux sends a boolean along the response channel
//A false value indicates that it was able to conclude definitively
//that n is composite (not prime)
func millerRabinAux(n big.Int, response chan bool) {
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
		response <- true
		return
	}

	s_minus_one := big.NewInt(0)
	s_minus_one = s_minus_one.Sub(s, big.NewInt(1))

	for i := big.NewInt(0); i.Cmp(s_minus_one) == -1; i.Add(i, big.NewInt(1)) {
		x = Exp(x, *big.NewInt(2), n)
		if x.Cmp(big.NewInt(1)) == 0 {
			response <- false
			return
		}

		tmp := big.NewInt(0)
		if x.Cmp(tmp.Sub(&n, big.NewInt(1))) == 0 {
			response <- true
			return
		}
	}
	response <- false
	return
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

func FindPrimeAndGenerator(n int64, certainty int) (big.Int, big.Int) {
	p := RandomNBitSafePrime(n, certainty)
	q := big.NewInt(0)
	q = q.Sub(&p, big.NewInt(1))
	q = q.Div(q, big.NewInt(2))
	nBig := big.NewInt(n)
	gCandidate := big.NewInt(1)
	for ; gCandidate.Cmp(&p) == -1; gCandidate = gCandidate.Add(gCandidate, big.NewInt(1)) {

		if e := Exp(*gCandidate, *big.NewInt(2), *nBig); e.Cmp(gCandidate) == 0 {
			continue
		}
		if e := Exp(*gCandidate, *q, *nBig); e.Cmp(gCandidate) == 0 {
			continue
		}
		break
	}
	return p, *gCandidate
}

func RSA(x *big.Int, bitlength int64, certainty int) (encoded, e, n, d *big.Int) {

	log.Print("Starting RSA")
	p := big.NewInt(0)
	q := big.NewInt(0)
	*p = RandomNBitPrime(bitlength/2, certainty)
	*q = RandomNBitPrime(bitlength/2, certainty)
	log.Printf("Generated primes %s and %s", p.String(), q.String())

	//phi = (p-1)(q-1)
	phi := big.NewInt(0)
	phi.Sub(p, big.NewInt(1))
	q_1 := big.NewInt(0)
	q_1 = q_1.Sub(q, big.NewInt(1))
	phi.Mul(phi, q_1)

	log.Printf("phi is %s", phi.String())

	n = big.NewInt(0)
	n = n.Mul(p, q)

	log.Printf("n is %s", n.String())
	//Generate e (this can be constant)
	e_int := int64(math.Floor(math.Pow(2, 16) + 1))
	e = big.NewInt(e_int)

	//Verify that phi is greater than (2^16 + 1)
	//Verify that gcd(e, phi) = 1
	if !(phi.Cmp(e) == 1) || gcd(e, phi).Cmp(big.NewInt(1)) != 0 {
		//We need to regenerate e
		//TODO fix this
		panic(fmt.Errorf("wrong value for e"))
	}

	encoded = big.NewInt(0)
	*encoded = Exp(*x, *e, *n)
	log.Printf("econded is %s", encoded.String())

	d = invert(e, phi)
	log.Printf("e is %s", e.String())

	return encoded, e, n, d
}

func RSA_Trapdoor(encoded, n, d *big.Int) (message *big.Int) {
	message = big.NewInt(0)
	*message = Exp(*encoded, *d, *n)
	return
}

func main() {
}
