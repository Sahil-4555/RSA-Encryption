package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// gcdExtended computes the extended Euclidean algorithm to find the greatest common divisor and the Bezout's coefficients.
func gcdExtended(a, b *big.Int, x, y *big.Int) *big.Int {
	zero := big.NewInt(0)

	if b.Cmp(zero) == 0 {
		x.Set(big.NewInt(1))
		y.Set(big.NewInt(0))
		return a
	}

	var x1, y1 big.Int
	gcd := gcdExtended(b, new(big.Int).Mod(a, b), &x1, &y1)

	x.Set(&y1)
	y.Set(new(big.Int).Sub(&x1, new(big.Int).Div(a, b).Mul(&y1, b)))

	return gcd
}

// modularExponentiation performs modular exponentiation using the repeated squaring method.
func modularExponentiation(a, b, m *big.Int) *big.Int {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	res := new(big.Int).Set(one) // Initialize result

	a.Mod(a, m) // Update a if it is more than or equal to m

	if a.Cmp(zero) == 0 {
		return zero // In case a is divisible by m
	}

	for b.Cmp(zero) > 0 {
		// If b is odd, multiply res with a and take modulo m
		if new(big.Int).And(b, one).Cmp(one) == 0 {
			res.Mul(res, a).Mod(res, m)
		}

		// b must be even now
		b.Rsh(b, 1)           // b = b/2 or b = b >> 1
		a.Mul(a, a).Mod(a, m) // a = (a^2) % m
	}

	return res
}

// modInverse calculates the modular inverse of num modulo modulus.
func modInverse(num, modulus *big.Int) *big.Int {
	gcd := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	// Compute the greatest common divisor (gcd) of num and modulus
	gcd.GCD(x, y, num, modulus)

	// If the gcd is not 1, the modular inverse doesn't exist
	if gcd.Cmp(big.NewInt(1)) != 0 {
		return big.NewInt(0)
	}

	// Make sure x is positive
	if x.Sign() < 0 {
		x.Add(x, modulus)
	}

	return x
}

// generateRandomCoprime generates random e value that is coprime with phi and falls within the specified range of greater than or equal to 2 and less than phi
func generateRandomCoprime(phi *big.Int) *big.Int {
	two := big.NewInt(2)
	one := big.NewInt(1)

	for {
		e, err := rand.Int(rand.Reader, phi)
		if err != nil {
			return big.NewInt(0)
		}

		// Check if e is coprime with phi and satisfies the condition
		if gcdExtended(e, phi, new(big.Int), new(big.Int)).Cmp(one) == 0 &&
			e.Cmp(two) >= 0 && e.Cmp(phi) < 0 {
			return e
		}
	}
}

// millerRabinTest performs the Miller-Rabin primality test on the given number n using k iterations.
func millerRabinTest(n *big.Int, k int) bool {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)
	four := big.NewInt(4)

	if n.Cmp(one) <= 0 || n.Cmp(four) == 0 {
		return false
	}
	if n.Cmp(three) <= 0 {
		return true
	}

	// Find r such that n = 2^d * r + 1 for some r >= 1
	d := new(big.Int).Sub(n, one)
	for new(big.Int).Mod(d, two).Cmp(zero) == 0 {
		d.Rsh(d, 1)
	}

	// Perform the Miller-Rabin test for 'k' iterations
	for i := 0; i < k; i++ {
		// Generate a random witness 'a' in the range [2, n-2] -> 1 < a < n-1
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, four))
		if err != nil {
			panic(err)
		}
		a.Add(a, two)

		// Compute a^(d % n)
		x := modularExponentiation(a, d, n)

		if x.Cmp(one) == 0 || x.Cmp(new(big.Int).Sub(n, one)) == 0 {
			continue
		}

		// Perform the Miller-Rabin test for 'd' iterations
		isPrime := false
		for r := new(big.Int).Set(one); r.Cmp(d) < 0; r.Lsh(r, 1) {
			x = modularExponentiation(x, two, n)

			if x.Cmp(one) == 0 {
				return false
			}
			if x.Cmp(new(big.Int).Sub(n, one)) == 0 {
				isPrime = true
				break
			}
		}

		if !isPrime {
			return false
		}
	}

	return true
}

// generateRandomPrime generates a random prime number with the specified number of bits.
func generateRandomPrime(numBits int) *big.Int {
	if numBits <= 0 {
		return big.NewInt(0)
	}

	minVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(numBits-1)), nil)
	maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(numBits)), nil)
	maxVal.Sub(maxVal, big.NewInt(1))

	for {
		num, err := rand.Int(rand.Reader, new(big.Int).Sub(maxVal, minVal))
		if err != nil {
			panic(err)
		}
		num.Add(num, minVal)

		if millerRabinTest(num, 10) {
			return num
		}
	}
}

func solve() {

	// Enter The Plain Text
	var m int64
	fmt.Println("Enter Your Message:")
	fmt.Scanln(&m)

	// Number of bits for the prime numbers
	numBits := 64
	one := big.NewInt(1)

	// generates two random prime number of numBits bit
	p := generateRandomPrime(numBits)
	q := generateRandomPrime(numBits)
	for q.Cmp(p) == 0 {
		q = generateRandomPrime(numBits)
	}

	fmt.Println("The Two Random Prime Numbers are:")
	fmt.Println("p:", p)
	fmt.Println("q:", q)

	// n is called Modulus of encryption and decryption.
	n := new(big.Int).Mul(p, q)
	fmt.Println("n:", n)

	/*
		choose a number e less than n, such that n is relatively prime to (p-1) * (q-1).
		It means that e and (p-1) * (q-1) have no common factor except 1. Therefore, gcd(e, phi(n)) = 1
	*/

	phi := new(big.Int).Mul(new(big.Int).Sub(p, one), new(big.Int).Sub(q, one))
	fmt.Println("Phi:", phi)

	// e value is generated that is coprime with phi
	e := generateRandomCoprime(phi)
	fmt.Println("e:", e)

	/*
		Then the public key is {e, n}.
		A plain text message m is encrypted using public key {e, n} to get ciphertext C.
	*/

	fmt.Println("Public Key {e,n}: {", e, ",", n, "}")

	/*
		To decrypt the received ciphertext C, use your private key exponent d.
		The private key exponent d is calculated by the formula:
		d = (k * Phi + 1) / e
	*/

	d := new(big.Int)
	d = modInverse(e, phi)
	cipherText := modularExponentiation(big.NewInt(m), e, n)
	fmt.Println("CipherText:", cipherText)

	/*
		Then the private key is {d, n}.
		A cipher text C is decrypted using private key {d, n} to get plaintext M.
	*/

	fmt.Println("d:", d)
	fmt.Println("Private Key {d,n}: {", d, ",", n, "}")

	// Decrypt the ciphertext using private key exponent d
	decipherText := modularExponentiation(cipherText, d, n)
	fmt.Println("DecipherText:", decipherText)
}

func main() {
	solve()
}
