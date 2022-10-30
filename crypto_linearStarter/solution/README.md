Use known plaintext `fla` to recover first three states: `s1 = ct[0] ^ ord("f")`. Solve linear equations:

s2 = s1 * A + B mod m

s3 = s2 * A + B mod m

To make the challenge easier, the modulus was chosen in a way that it does not matter in this calculation. Therefore the following equations can be used to recover A and B:

A = (s2 - s3) / (s1 - s2)

B = s2 - (A * s1)

Then use A and B to compute the complete key used for the OTP and recover the flag.
