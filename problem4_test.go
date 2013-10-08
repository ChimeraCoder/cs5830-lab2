package main

import (
	"log"
	"math"
	"math/big"
	"testing"
)

func Test_Exponentiation(t *testing.T) {

	if e := Exp(*big.NewInt(3), *big.NewInt(16), *big.NewInt(3000)); e.Cmp(big.NewInt(2721)) != 0 {
		t.Error("exp failed")
		log.Print("first exp failed")
	} else {
		t.Log("first exp passed") // log some info if you want
		log.Print("first exp passed")
	}

	//if Exp(2, 10, 5) != 4 {
	if e := Exp(*big.NewInt(2), *big.NewInt(10), *big.NewInt(5)); e.Cmp(big.NewInt(4)) != 0 {
		t.Error("second exp failed")
		log.Print("second exp failed")
	} else {
		t.Log("second exp passed")
		log.Print("second exp passed")
	}
}

func Test_invert(t *testing.T) {
	//Check if inverting 101 and 102 yields -1, as it should
	log.Print("Testing inverting")
	if invert(big.NewInt(101), big.NewInt(102)).Cmp(big.NewInt(-1)) != 0 {
		t.Error("Invert failed")
		log.Print("Invert failed to invert 101 mod 102 to -1")
	} else {
		log.Print("Invert passed")
	}
}

func Test_MillerRabin(t *testing.T) {

	prime_tests := []int64{23251, 999331, 115249, 479001599}

	for _, prime := range prime_tests {
		prime_big := big.NewInt(prime)
		if !MillerRabin(*prime_big, 10) {
			t.Errorf("MillerRabin failed to detect primality of %s", prime_big.String())
		} else {
			log.Printf("Succeeded on prime %s", prime_big.String())
		}

		composite := prime + 2
		composite_big := big.NewInt(composite)
		if MillerRabin(*composite_big, 10) {
			t.Errorf("MillerRabin claimed that composite number %d is prime", composite)
		} else {
			log.Printf("Succeeded on composite number %d", composite)
		}
	}
}

func TestRandomNBitNumber(t *testing.T) {

	number := RandomNBitNumber(62)
	lower := big.NewInt(0)
	upper := big.NewInt(0)
	lower, s1 := lower.SetString("4611686018427387904", 10)
	upper, s2 := upper.SetString("9223372036854775808", 10)

	if !s1 || !s2 {
		t.Errorf("Error creating bounds for test")
	}

	if lower.Cmp(&number) > 0 || upper.Cmp(&number) != 1 {
		t.Errorf("Attempting to generate a %d-bit number yielded %s", 62, number.String())
	}
}

func Test_RandomNBitSafePrime(t *testing.T) {
	//bits := int64(80)
	bits := int64(20)
	number := RandomNBitSafePrime(bits, 10)

	//If number is a safe prime, "other_prime" should be prime too

	other_prime := big.NewInt(0)
	other_prime = other_prime.Sub(&number, big.NewInt(1))
	other_prime = other_prime.Div(other_prime, big.NewInt(2))

	//This is equivalent to
	//other_prime := (number - 1) / 2

	if !MillerRabin(*other_prime, 20) {
		t.Errorf("Error - %s is not a safe prime, as %s is composite", number.String(), other_prime.String())
	} else {
		log.Printf("Succeded in generating safe prime %s", number.String())
	}
}

func Test_FindGenerator(t *testing.T) {
	bits := int64(6)
	certainty := 10
	p1, g1 := FindPrimeAndGenerator(bits, certainty)
	p := &p1
	g := &g1

	//g^2 and g^((p-1)/2) should not yield g again (mod p)

	error_found := false

	result := big.NewInt(0)
	result.Exp(g, big.NewInt(2), p)
	if result.Cmp(g) == 0 {
		t.Errorf("Falsely identified %s as a generator for Zp* for p = %s", g.String(), p.String())
		error_found = true
	}

	//q = (p-1)/2
	q := big.NewInt(0)
	q = q.Sub(p, big.NewInt(1))
	q = q.Div(q, big.NewInt(2))

	result = big.NewInt(0)
	result.Exp(g, q, p)
	if result.Cmp(g) == 0 {
		t.Errorf("Falsely identified %s as a generator for Zp* for p = %s", g.String(), p.String())
		error_found = true
	}

	if !error_found {
		log.Printf("Successfully identified generator %s for prime %s", g.String(), p.String())
	}

}

func Test_RSA(t *testing.T) {

	certainty := 40

	message := big.NewInt(0)

	//Message should be a random integer on the range [0, 2^20)
	message = big.NewInt(0).Rand(r, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(20), nil))
	log.Printf("Message is %s", message.String())
	bitLength := int64(math.Floor(math.Pow(2, 8)))
	//Together, e and n form the public key
	//Together, n and d form the public key
	//d = 1/e mod phi(n)

	log.Print("encoding")

	bitLength = 128

	encoded, e, n, d := RSA(message, bitLength, certainty)
	log.Print("encoded")
	log.Printf("d is %s", d.String())

	result := big.NewInt(0)
	if result.Exp(message, e, n).Cmp(encoded) != 0 {
		t.Errorf("The encoded message is not the same as x^e, mod N")
	}

	//The trapdoor does not need e
	decrypted := RSA_Trapdoor(encoded, n, d)

	if decrypted.Cmp(message) != 0 {
		t.Errorf("Encrypting with RSA and then decrypting message with RSA trapdoor did not yield the same result")
	} else {
		t.Log("Successfully encrypted message with RSA and decrypted with the RSA trapdoor")
	}

	log.Print("Encoding test message using various key bitlengths.")
	m := big.NewInt(0)
	m, re := m.SetString("22405534230753963835153736737", 10)
	if !re {
			panic("Failed to scan string literal into a BigInt")
	}
	log.Printf("Message: %s", m.String())
	log.Print("(\"Hello World!\"'s 96-bit binary representation, in decimal)")


	for _, i := range []int{128, 256, 512} {
		log.Print("----------")
		log.Printf("%d BITS", i)
		log.Printf("Finding %d-bit safe prime...", i)
		g := big.NewInt(0)
		p := big.NewInt(0)
		*p, *g = FindPrimeAndGenerator(int64(i), 4)
		log.Printf("p: %s", p.String())
		log.Print("Discrete log OWF:")
		log.Printf("g: %s", g.String())
		enc := big.NewInt(0)
		*enc = Exp(*g, *m, *p)
		log.Printf("Enc(message): %s", enc.String())
		log.Print("\n")
		log.Print("RSA OWF:")
		n := big.NewInt(0)
		d := big.NewInt(0)
		e := big.NewInt(0)
		enc, e, n, d = RSA(m, int64(i), 4)
		log.Printf("Enc(message): %s", enc.String())
		log.Printf("n: %s", n.String())
		log.Printf("e: %s", e.String())
		log.Printf("d: %s", d.String())
		dec := big.NewInt(0)
		*dec = Exp(*enc, *d, *n)
		if m.Cmp(dec) == 0 {
			log.Printf("Dec(message): %s", dec.String())
			log.Print("... which == the message.");
		} else {
			panic("RSA decoding produced value different from message")
		}
	}

}
