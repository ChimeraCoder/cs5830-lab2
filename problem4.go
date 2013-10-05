package main

import (
	"math"
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
    //check if pow == 0
	if pow.Cmp(big.NewInt(0)) == 0 {
		return *big.NewInt(1)
	}

    //log.Print("Done comapring")
	//We need to square a "square_iterations" times
	//then multiply by a^(pow-(2^square_iterations))
    //Without math.Big, this would be equivalent to the following:
	//square_iterations := int(math.Floor(math.Log2(float64(pow))))
    square_iterations := pow.BitLen()


	tmp := big.NewInt(1)
	for i := 0; i < square_iterations; i++ {
        //tmp = tmp * tmp % n
		tmp = tmp.Mul(tmp, tmp)
        
        tmp = tmp.Mod(tmp, &n)
	}
    //log.Print("Done iterating powers of 2")

	//Now, tmp = a^(2^square_iterations)

	//TODO make this more performant


    //how many multiplications have already been performed
    multiplications_performed := big.NewInt(0)
    multiplications_performed = multiplications_performed.Exp(big.NewInt(2), big.NewInt(int64(square_iterations)), nil)
    
    new_pow := big.NewInt(0)
    new_pow = new_pow.Sub(&pow, multiplications_performed)


    //Without math/big, this would be the following
    //return tmp * Exp(a, new_pow, n) % n
    result := big.NewInt(0)
    //log.Printf("About to do recursive call on %v %v %v", a, *new_pow, n) 
    *result = Exp(a, *new_pow, n)
    result = result.Mul(tmp, result)
    result = result.Mod(result, &n)
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

func MillerRabin(n int64, numTests int) bool {
    for i := 0; i < numTests; i++ {
        if (MillerRabinAux(n) == false) {
            return false
        }
    }
    return true
}

func MillerRabinAux(n int64) bool {
    d := n-1
    s := 0
    for ; d%2 == 0; {
        d = d / 2
        s++
    }

    a := r.Int63n(n-3) + 2 // returns vals between [2, n-2] incl.
    x := Exp(*big.NewInt(a), *big.NewInt(d), *big.NewInt(n))


    //Equivalent to
    //if (x == 1) || (x == n-1) {
    if (x.Cmp(big.NewInt(1)) == 0) || (x.Cmp(big.NewInt(n-1)) == 0) {
        return true
    }

    for i := 0; i < s-1; i++ {
        x = Exp(x, *big.NewInt(2), *big.NewInt(n))
        if (x.Cmp(big.NewInt(1)) == 0) {
            return false
        }
        if (x.Cmp(big.NewInt(n-1)) == 0) {
            return true
        }
    }
    return false
}


//RandomNBitNumber returns a random number with the specified number of bits
func RandomNBitNumber(n int) int64{
    lower_bound := int64(math.Pow(2, float64(n))) //Inclusive
    upper_bound := int64(math.Pow(2, float64(n+1))) //Exclusive

    return r.Int63n(upper_bound - lower_bound) + lower_bound
}

//RandomNBitPrime returns random prime numbers of the specified size
func RandomNBitPrime(n int, certainty int) int64{
    for {
        n := RandomNBitNumber(n)
        if MillerRabin(n, certainty){
            return n
        }
    }
}

func main() {
}
