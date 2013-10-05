package main

import (
	"testing"
    "log"
    "math/big"
)

func Test_Exponentiation(t *testing.T) {

	if e := Exp(*big.NewInt(3), *big.NewInt(16), *big.NewInt(3000)); e.Cmp(big.NewInt(2721)) != 0 {
		t.Error("exp failed")
	} else {
		t.Log("first exp passed") // log some info if you want
	}
	//if Exp(2, 10, 5) != 4 {
	if e := Exp(*big.NewInt(2), *big.NewInt(10), *big.NewInt(5)); e.Cmp(big.NewInt(4)) != 0 {
		t.Error("second exp failed")
	}
}

func Test_invert(t *testing.T) {
	if invert(101, 102) != -1 {
		t.Error("Invert failed")
	}
}

func Test_MillerRabin(t *testing.T) {

    prime_tests := []int64{23251, 999331, 115249, 479001599}

    for _, prime := range prime_tests{
        if !MillerRabin(prime, 10){
            t.Errorf("MillerRabin failed to detect primality of %d", prime)
        } else {
            log.Printf("Succeeded on prime %d", prime)
        }

        composite := prime +2 
        if MillerRabin(composite, 10) {
            t.Errorf("MillerRabin claimed that composite number %d is prime", composite)
        } else {
            log.Printf("Succeeded on composite number %d", composite)
        }
    }
}


/**
func TestRandomNBitNumber(t *testing.T){
    
    bits := 62
    number := RandomNBitNumber(bits)
    if !(4611686018427387904 <= number &&  number < 9223372036854775808) {
        t.Errorf("Attempting to generate a %d-bit number yielded %d", bits, number)
    }
}
*/
