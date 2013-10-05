package main

import (
	"testing"
    "log"
)

func Test_Exponentiation(t *testing.T) {
	if Exp(3, 16, 3000) != 2721 {
		t.Error("exp failed")
	} else {
		t.Log("first exp passed") // log some info if you want
	}
	if Exp(2, 10, 5) != 4 {
		t.Error("second exp failed")
	}
}

func Test_invert(t *testing.T) {
	if invert(101, 102) != -1 {
		t.Error("Invert failed")
	}
}

func Test_MillerRabin(t *testing.T) {

    
    prime_tests := []int{23251, 999331, 115249, 479001599}


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

