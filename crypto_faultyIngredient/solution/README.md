DISCLAIMER: The code of `solve.go` is really not good, but it does the trick :) 

At first the $MC([\Delta,0,0,0]), MC([0,\Delta,0,0]), MC([0,0,\Delta,0])$ and $MC([0,0,0,\Delta])$ values are precomputed using an own implementation of the AES-MixColumn function. The 4 byte results are stored in `uint32`.

Then the fault attack repeats for every one of the four columns:

1. Attack $k_1$ and $k_2$:
   + Compute every guess as 16-bit integer.
   + For every guess compute the SubBytes Input for the non-faulty ciphertext and the faulty ciphertext to compute the difference using XOR
   + Check whether the difference (16-bit value) is contained in one of the precomputed deltas: `(delta >> 16) == diff)`
     + If included: Continue
     + else: Discard from guesses
2. Repeat same attack for $k_3$ using only the left over guess from step 1: `key_space_k12[uint32(k12)]<<8 ^ uint32(k3)`
3. Finally add $k_4$ in the same way and run attack on using the full 32-bit precomputed deltas and guesses.
   + Repeat this step until only one guess remains.

Then reverse the AES-128 Key-Scheduling and use main key to decrypt the flag.
