#include <iostream>
#include <bits/stdc++.h>
#include <cstdint>
#include <random>
#include <cmath>
using namespace std;
#define FAST                     \
    ios::sync_with_stdio(false); \
    cin.tie(NULL);               \
    cout.tie(NULL);

// Function to calculate the greatest common divisor (GCD) of two numbers using extended Euclidean algorithm
uint64_t gcdExtended(uint64_t a, uint64_t b, uint64_t *x, uint64_t *y)
{
    if (b == 0)
    {
        *x = 1;
        *y = 0;
        return a;
    }

    uint64_t x1, y1;
    uint64_t gcd = gcdExtended(b, a % b, &x1, &y1);

    *x = y1;
    *y = x1 - (a / b) * y1;

    return gcd;
}

// Function to calculate (a^b) % m using binary exponentiation
uint64_t modularExponentiation(uint64_t a, uint64_t b, uint64_t m)
{
    uint64_t res = 1; // Initialize result

    a = a % m; // Update a if it is more than or equal to m

    if (a == 0)
        return 0; // In case a is divisible by m

    while (b > 0)
    {
        // If b is odd, multiply res with a and take modulo m
        if (b & 1)
            res = ((res % m) * (a % m)) % m;

        // b must be even now
        b = b >> 1; // b = b/2
        a = ((a % m) * (a % m)) % m;
    }
    return res;
}

// Function to calculate the modular multiplicative inverse of 'e' modulo 'tot_n'
uint64_t modInverse(uint64_t e, uint64_t tot_n)
{
    int64_t r1 = tot_n, r2 = e;
    int64_t t1 = 0, t2 = 1;
    int64_t q, r, t;

    while (r2 != 0)
    {
        q = r1 / r2;
        r = r1 - q * r2;
        r1 = r2;
        r2 = r;
        t = t1 - q * t2;
        t1 = t2;
        t2 = t;
    }

    // Make sure the result is positive
    uint64_t inverse = (t1 < 0) ? (t1 + tot_n) : t1;
    return inverse;
}

// generateRandomCoprime generates random e value that is coprime with phi and falls within the specified range of greater than or equal to 2 and less than phi
uint64_t generateRandomCoprime(uint64_t phi)
{
    while (true)
    {
        srand(time(NULL));
        uint64_t e = rand() % phi;
        uint64_t x, y;
        // e must be co-prime and smaller than phi
        if (gcdExtended(e, phi, &x, &y) == 1 && e >= 2 && e < phi)
            return e;
    }
}

// Function to perform the Miller-Rabin primality test
bool millerRabinTest(uint64_t n, int k)
{
    if (n <= 1 || n == 4)
        return false;
    if (n <= 3)
        return true;

    // Find r such that n = 2^d * r + 1 for some r >= 1
    uint64_t d = n - 1;
    while (d % 2 == 0)
        d >>= 1;

    // Perform the Miller-Rabin test for 'k' iterations
    for (int i = 0; i < k; ++i)
    {
        // Generate a random witness 'a' in the range [2, n-2]
        random_device rd;
        mt19937_64 gen(rd());
        uniform_int_distribution<uint64_t> dis(2, n - 2);
        uint64_t a = dis(gen);

        // Compute a^(d % n)
        uint64_t x = modularExponentiation(a, d, n);

        if (x == 1 || x == n - 1)
            continue;

        // Perform the Miller-Rabin test for 'd' iterations
        bool isPrime = false;
        for (uint64_t r = 1; r < d; r <<= 1)
        {
            x = modularExponentiation(x, 2, n);

            if (x == 1)
                return false;
            if (x == n - 1)
            {
                isPrime = true;
                break;
            }
        }

        if (!isPrime)
            return false;
    }

    return true;
}

// Function to generate a random prime number with the given number of bits
uint64_t generateRandomPrime(int numBits)
{
    if (numBits <= 0)
        return 0;

    srand(time(NULL));
    mt19937_64 gen(time(nullptr));

    uint64_t minVal = pow(2, numBits - 1);
    uint64_t maxVal = pow(2, numBits) - 1;

    uniform_int_distribution<uint64_t> dis;

    uint64_t num = 0;
    do
    {
        num = dis(gen, decltype(dis)::param_type(minVal, maxVal));
    } while (!millerRabinTest(num, 10));

    return num;
}

void solve()
{
    int numBits = 16; // Number of bits for the prime numbers

    // Enter The Plain Text
    string m;
    cout << "Enter Message: " << endl;
    getline(cin, m);

    uint64_t p = generateRandomPrime(numBits);
    uint64_t q = p;
    while (q == p)
    {
        q = generateRandomPrime(numBits);
    }

    cout << "The Two Random Numbers are: " << endl;
    cout << "p: " << p << endl;
    cout << "q: " << q << endl;

    // n is called Modulus of encryption and decryption.
    uint64_t n = p * q;
    cout << "n: " << n << endl;

    /*
        choose a number e less than n, such that n is relatively
        prime to (p-1) * (q-1). It means that e and (p-1) * (q-1)
        have no common factor except 1. Therefore, gcd(e, phi(n)) = 1
    */

    uint64_t phi = (p - 1) * (q - 1);

    // Generate e value that is coprime with phi
    uint64_t e = generateRandomCoprime(phi);
    cout << "Phi: " << phi << endl;

    cout << "e: " << e << endl;

    cout << "Public Key: {" << e << ", " << n << "}" << endl;

    /*
        Then the public key is {e, n}. A plain text message m is encrypted
        using public key {e, n} to get ciphertext C.
    */
    vector<uint64_t> CipherText;
    for (int i = 0; i < m.size(); i++)
    {
        uint64_t ct = modularExponentiation(m[i] - 0, e, n);
        CipherText.push_back(ct);
    }
    cout << "CipherText: ";
    for (auto it : CipherText)
    {
        cout << char(it % 128);
    }
    cout << endl;

    /*
        To decrypt the received ciphertext C, use your private key exponent d.
        The private key exponent d is calculated by the formula:
        d = (k * Phi + 1) / e
    */
    uint64_t d = modInverse(e, phi);
    cout << "d: " << d << endl;
    cout << "Private Key: {" << d << ", " << n << "}" << endl;

    // Decrypt the ciphertext using private key exponent d
    string DecipherText;
    for (auto it : CipherText)
    {
        uint64_t dt = modularExponentiation(it, d, n) + 0;
        DecipherText.push_back(dt);
    }
    cout << "DecipherText: " << DecipherText << endl;
}

int main()
{

    FAST;
    solve();

    return 0;
}
