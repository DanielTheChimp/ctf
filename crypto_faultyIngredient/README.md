# Faulty Ingredient
- Category: **Crypto**
- Difficulty: **Medium / Hard** (Depends heavily on whether you have done such an attack before)
- CTF: **Hack.lu 2022**

## Description
My boss said, that we should better encrypt the secret ingredient of our delicious sauce. But when I used my own AES, I noticed that sometimes, I get totally wrong values. Since I had a great crypto class at my evening school, I managed to nail the error down to the **input of the Mix Columns in Round 9**. I can't say how, but when the error occurs, **exactly one random byte in every column** of the state changes. I ran some tests, solely for myself, but my dumb intern posted them on fluxstagram when making one of his stupid selfies. Well anyway, i hope nobody can use this to recover the **main key** and decrypt our secret ingredient!

## Setup
None
