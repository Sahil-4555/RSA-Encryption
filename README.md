#  RSA algorithm (Rivest-Shamir-Adleman)

This code implements the RSA encryption and decryption algorithm. It generates random prime numbers, calculates the public and private keys, encrypts and decrypts messages using these keys.

## Functionality

The code performs the following steps:

1. Generates two random prime numbers, `p` and `q`, of the specified number of bits.

2. Calculates the modulus `n` as the product of `p` and `q`.

3. Calculates the Euler's totient function `phi` as `(p-1) * (q-1)`.

4. Finds the optimal value of `e` that is relatively prime to `phi`.

5. Displays the generated prime numbers (`p` and `q`), modulus (`n`), Euler's totient function (`phi`), and the public key (`e`, `n`).

6. Encrypts the input plain text message using the public key to obtain the ciphertext.

7. Displays the encrypted ciphertext.

8. Calculates the private key `d` as the modular multiplicative inverse of `e` modulo `phi`.

9. Decrypts the ciphertext using the private key to obtain the original plain text message.

10. Displays the decrypted plain text message.

## Go Implementation

Run The `main.go` file Using Command:
```
go run main.go
````
### Output
<div align="center">

![RSA_GO_IMPLEMENTATION](https://www.linkpicture.com/q/RSA_GO_3.png)

</div>

## C++ Implementation

Run The `RSA.cpp` file Using Command:
```
g++ RSA.cpp -o rsa
````
Once the compilation process is successful, you can execute the compiled program. To run the program, use the following command:
```
./rsa
```
### Output

![RSA_C++_IMPLEMENTATION](https://www.linkpicture.com/q/RSA_C.png)

## Customization

- You can modify the number of bits for the generated prime numbers by changing the `numBits` variable in the `solve()` function.

## Acknowledgments

- The code is based on the RSA algorithm for encryption and decryption.
- The Miller-Rabin primality test is used to generate random prime numbers.
- The modular exponentiation and extended Euclidean algorithm functions are adapted from standard algorithms.
