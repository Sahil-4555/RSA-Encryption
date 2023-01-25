#include <bits/stdc++.h>
using namespace std;
#define ll long long
#define v vector<ll>
const ll N = 1e5 + 1;
#define FAST                     \
    ios::sync_with_stdio(false); \
    cin.tie(NULL);               \
    cout.tie(NULL);

vector<ll> primeConatiner;

// filling prime number for random picking in container

void FillPrimeNumber()
{
    vector<bool> sieve(100001, true);
    sieve[0] = sieve[1] = false;

    // marking false which are not prime number.

    for (ll i = 2; i * i <= 100000; i++)
        if (sieve[i] == true)
            for (ll j = i * i; j <= 100000; j += i)
                sieve[j] = false;

    /*
    pushing number in container which are true bcz they are
    prime number.
    */

    // ofstream outdata;
    // outdata.open("primes.txt");
    for (ll i = 0; i < sieve.size(); i++)
        if (sieve[i])
            primeConatiner.push_back(i);
    // outdata << i << " ";

    // outdata.close();
}

ll pickRandomPrime()
{
    srand(time(NULL));
    // choosing random index from primeContainer.
    ll k = rand() % primeConatiner.size();

    // marking the starting position as first index
    auto itr = primeConatiner.begin() + k;

    // looping until we get to picked index.
    // while (k--)
    // {
    //     itr++;
    // }

    /*
    erasing the picked number so no duplicate number been picked
    second time.
    */

    ll data = primeConatiner[k];
    // auto index = find(primeConatiner.begin(), primeConatiner.end(), data);
    primeConatiner.erase(itr);

    // returning the picked prime   number.
    return data;
}

// Converting Decimal to Binary

vector<ll> ConvertNumToBinary(ll num)
{
    // string to store binary number
    vector<ll> binaryNum(32);

    // counter for binary array
    ll i = 0;

    while (num > 0)
    {
        // storing remainder in binary array
        binaryNum[i] = num % 2;
        num = num / 2;
        i++;
    }

    // erasing extra zeros in 32-bit binary array
    binaryNum.erase(binaryNum.begin() + i, binaryNum.end());

    // storing binary array in reverse order
    reverse(binaryNum.begin(), binaryNum.end());

    // returning the binary array of number num
    return binaryNum;
}

// Square and Multiply Algorithm OR Fast Modulor Exponentiation
// (pow(a,b) % c) OR (a^b mod c)

ll calculate(ll a, ll b, ll c)
{
    vector<ll> binaryOfB = ConvertNumToBinary(b);

    ll i = 1;
    ll prev = a;
    while (i < binaryOfB.size())
    {
        ll tmp = pow(prev, 2);
        if (binaryOfB[i])
            prev = ((tmp % c) * a) % c;
        else
            prev = tmp % c;
        i++;
    }

    return prev;
}

ll gcdExtended(ll a, ll b, ll *x, ll *y)
{
    // Base Case
    if (a == 0)
    {
        *x = 0;
        *y = 1;
        return b;
    }

    ll x1, y1; // To store results of recursive call
    ll gcd = gcdExtended(b % a, a, &x1, &y1);

    // Update x and y using results of
    // recursive call
    *x = y1 - (b / a) * x1;
    *y = x1;

    return gcd;
}

void solve()
{
    // Fill the prime number in container.

    FillPrimeNumber();

    // Select Two Random Large Prime Numbers.

    ll p = pickRandomPrime();
    ll q = pickRandomPrime();

    cout << "The Two Random Number are: " << endl;
    cout << "p: " << p << endl;
    cout << "q: " << q << endl;

    // Enter The Plain Text

    ll m;
    cin >> m;

    // n is called Modulous of encrption and decryption.

    ll n = p * q;

    /*
        choose a number e less then n,such that n is relatively
        prime to (p-1) (q-1). it means that e and (p-1)*(q-1)
        have no common factor except 1 therefore gcd(e,phi(n)) = 1
    */

    ll e = 2;
    ll phi = (p - 1) * (q - 1);
    cout << "Phi: " << phi << endl;
    while (e < phi)
    {
        ll x, y;
        // e must be co-prime and smaller than phi
        if (gcdExtended(e, phi, &x, &y) == 1)
            break;
        else
            e++;
    }

    cout << "e: " << e << endl;

    cout << "Public Key: {" << e << "," << n << "}" << endl;

    /*
    then the public key is {e,n}. a plain text message m is encrypted
    using public key {e,n} to get ciphertext C
    */

    // c = pow(m,e) % n OR m^e mod n

    ll c = calculate(m, e, n);

    cout << "Ciphet Text: " << c << endl;

    // To determine the private key

    ll k = 0;
    while (floor(double(double(1 + (k * phi)) / double(e))) != ceil(double(double(1 + (k * phi)) / double(e))))
    {
        k++;
    }

    cout << "k: " << k << endl;

    ll d = (1 + (k * phi)) / e;

    cout << "Private Key: {" << d << "," << n << "}" << endl;

    /*
    the private key is {d,n}. a ciphertext message c is decrypted using
    private key {d,n}. to calculate plain text m from ciphertext
    */

    // m = pow(c,d) % n OR c^d mod n

    m = calculate(c, d, n);

    cout << "Plain Text: " << m << endl;
}

int main()
{

    FAST;
#ifndef ONLINE_JUDGE
    freopen("input", "r", stdin);
    freopen("output", "w", stdout);
#endif

    solve();

    return 0;
}