package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

// gcdExtended computes the extended Euclidean algorithm to find the greatest common divisor and the Bezout's coefficients.
func gcdExtended(a, b *big.Int, x, y *big.Int) *big.Int {
	zero := big.NewInt(0)

	// If b is zero
	if b.Cmp(zero) == 0 {
		// Set x to 1
		x.Set(big.NewInt(1))
		// Set y to 0
		y.Set(big.NewInt(0))
		// Return a
		return a
	}

	var x1, y1 big.Int
	// Recursive call with updated values
	gcd := gcdExtended(b, new(big.Int).Mod(a, b), &x1, &y1)

	// Set x to the value of y1
	x.Set(&y1)

	// Set y to the subtraction of x1 and the product of (a/b) and y1
	y.Set(new(big.Int).Sub(&x1, new(big.Int).Div(a, b).Mul(&y1, b)))

	// Return the gcd
	return gcd
}

// modularExponentiation performs modular exponentiation using the repeated squaring method.
func modularExponentiation(a, b, m *big.Int) *big.Int {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	// Initialize result to 1
	res := new(big.Int).Set(one)

	// Update a if it is more than or equal to m
	a.Mod(a, m)

	if a.Cmp(zero) == 0 {
		// If a is divisible by m, return 0
		return zero
	}

	// While b is greater than 0
	for b.Cmp(zero) > 0 {
		// If b is odd, multiply res with a and take modulo m
		if new(big.Int).And(b, one).Cmp(one) == 0 { // Check if b is odd
			// Multiply res with a and take modulo m
			res.Mul(res, a).Mod(res, m)
		}

		// b must be even now
		// Right shift b by 1 (equivalent to b/2 or b = b >> 1)
		b.Rsh(b, 1)
		// a = (a^2) % m
		// Square a and take modulo m
		a.Mul(a, a).Mod(a, m)
	}

	// Return the result
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
		// Return 0 if modular inverse doesn't exist
		return big.NewInt(0)
	}

	// Make sure x is positive
	if x.Sign() < 0 { // If x is negative
		// Add modulus to x
		x.Add(x, modulus)
	}

	// Return the modular inverse
	return x
}

// generateRandomCoprime generates random e value that is coprime with phi and falls within the specified range of greater than or equal to 2 and less than phi
func generateRandomCoprime(phi *big.Int) *big.Int {
	two := big.NewInt(2)
	one := big.NewInt(1)

	for {
		// Generate a random integer e between 0 and phi-1
		e, err := rand.Int(rand.Reader, phi)
		if err != nil {
			// Return 0 if there is an error
			return big.NewInt(0)
		}

		// Check if e is coprime with phi and satisfies the condition
		if gcdExtended(e, phi, new(big.Int), new(big.Int)).Cmp(one) == 0 && // If e and phi are coprime
			e.Cmp(two) >= 0 && e.Cmp(phi) < 0 { // If e is greater than or equal to 2 and less than phi
			return e // Return e
		}
	}
}

// millerRabinTest performs the Miller-Rabin primality test to check if a number is probably prime.
func millerRabinTest(n *big.Int, k int) bool {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)

	if n.Cmp(two) == 0 || n.Cmp(three) == 0 { // If n is 2 or 3, return true
		return true
	}

	if n.Cmp(two) < 0 || new(big.Int).And(n, one).Cmp(zero) == 0 { // If n is less than 2 or even, return false
		return false
	}

	// Express n-1 as 2^r * d, where d is an odd number
	r := 0
	d := new(big.Int).Sub(n, one)
	for d.Bit(0) == 0 { // While the least significant bit of d is 0
		r++         // Increment r
		d.Rsh(d, 1) // Right shift d by 1 (equivalent to d/2)
	}

	for i := 0; i < k; i++ { // Repeat the test for k iterations
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, three)) // Generate a random number a between 2 and n-2
		if err != nil {
			return false // Return false if there is an error
		}
		a.Add(a, two) // Add 2 to a to make it fall within the range of [2, n-2]

		x := modularExponentiation(a, d, n) // Compute a^d modulo n

		if x.Cmp(one) == 0 || x.Cmp(new(big.Int).Sub(n, one)) == 0 { // If x is 1 or n-1
			continue // Continue to the next iteration
		}

		for j := 1; j < r; j++ { // Repeat r-1 times
			x.Mul(x, x).Mod(x, n) // Square x and take modulo n

			if x.Cmp(one) == 0 { // If x is 1, return false (composite)
				return false
			}

			if x.Cmp(new(big.Int).Sub(n, one)) == 0 { // If x is n-1
				break // Break the inner loop and continue to the next iteration
			}
		}

		if x.Cmp(new(big.Int).Sub(n, one)) != 0 { // If x is not n-1, return false (composite)
			return false
		}
	}

	return true // If the number passes all iterations, return true (probably prime)
}

// generateRandomPrime generates a random prime number with the specified bit length.
func generateRandomPrime(numBits int) *big.Int {
	// Return 0 if numBits is less than or equal to 0
	if numBits <= 0 {
		return big.NewInt(0)
	}

	// Calculate the minimum value as 2^(numBits-1)
	minVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(numBits-1)), nil)

	// Calculate the maximum value as 2^numBits
	maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(numBits)), nil)

	// Subtract 1 from the maximum value
	maxVal.Sub(maxVal, big.NewInt(1))

	for {
		// Generate a random number between minVal and maxVal
		num, err := rand.Int(rand.Reader, new(big.Int).Sub(maxVal, minVal))
		if err != nil {
			// Panic if there is an error generating the random number
			panic(err)
		}

		// Add minVal to the generated random number
		num.Add(num, minVal)

		// If the number passes the Miller-Rabin test with 20 iterations
		if millerRabinTest(num, 20) {
			// Return the prime number
			return num
		}
	}
}

func solve() {

	// Enter The Plain Text
	var inputText string
	fmt.Println("Enter Input:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputText = scanner.Text()

	// Number of bits for the prime numbers
	numBits := 64
	one := big.NewInt(1)

	// generates two random prime number of numBits bit
	p := generateRandomPrime(numBits)
	q := generateRandomPrime(numBits)
	for q.Cmp(p) == 0 {
		q = generateRandomPrime(numBits)
	}

	fmt.Println("The Two Random Prime Numbers are: ")
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

	// Generate e value that is coprime with phi
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

	// // Convert the individual characters of the string to their ASCII values
	plainText := []rune(inputText)

	encryptedText := make([]*big.Int, len(plainText))

	for i, char := range plainText {
		eValue := new(big.Int).Set(e)
		cipherText := modularExponentiation(big.NewInt(int64(char)), eValue, n)
		encryptedText[i] = cipherText
	}

	encryptedStr := ""
	for _, cipherText := range encryptedText {
		var tmp big.Int
		tmp.Mod(cipherText, big.NewInt(128)) // Modulo 128
		char := rune(tmp.Int64())            // Convert to ASCII value
		encryptedStr += string(char)
	}

	fmt.Println("Encrypted Text:", encryptedStr)

	/*
		Then the private key is {d, n}.
		A cipher text C is decrypted using private key {d, n} to get plaintext M.
	*/

	fmt.Println("d:", d)
	fmt.Println("Private Key {d,n}: {", d, ",", n, "}")

	// Decrypt the ciphertext using private key exponent d
	decryptedStr := ""

	for _, char := range encryptedText {
		dValue := new(big.Int).Set(d)
		decipherText := modularExponentiation(char, dValue, n)
		decryptedStr += string(rune(decipherText.Int64()))
	}

	fmt.Println("Decrypted Text:", (decryptedStr))
}

func main() {
	solve()
}
