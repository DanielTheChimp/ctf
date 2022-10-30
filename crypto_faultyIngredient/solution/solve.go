package main

import (
	"crypto/aes"
	"fmt"
	"os"
	"time"
	"encoding/json"
)

/*	 
* HELPER STUFF FOR GO
*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func copyMap(og map[uint32]uint32) map[uint32]uint32 {
	copy := make(map[uint32]uint32, len(og))
	for k, v := range og {
		copy[k] = v
	}
	return copy
}

func sliceToBytes(data []uint) []byte {
	res := make([]byte, 16)
	for i := range res {
		res[i] = byte(data[i])
	}
	return res
}

func sliceToHexStrings(data []uint) []string {
	res := make([]string, 16)
	for i := range res {
		res[i] = fmt.Sprintf("%x", data[i])
	}
	return res
}

/*	 
* AES STUFF
*/

func sbox() []uint {
	return []uint{
		0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
		0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
		0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
		0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
		0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
		0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
		0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
		0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
		0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
		0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
		0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
		0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
		0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
		0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
		0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
		0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16}
}

func inv_sbox() []uint {
	return []uint{
		0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb,
		0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb,
		0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e,
		0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25,
		0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92,
		0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84,
		0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06,
		0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b,
		0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73,
		0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e,
		0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b,
		0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4,
		0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f,
		0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef,
		0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61,
		0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d}
}

func shift_rows() []uint {
	return []uint{0, 5, 10, 15, 4, 9, 14, 3, 8, 13, 2, 7, 12, 1, 6, 11}
}

func inv_shift_rows() []uint {
	return []uint{0, 13, 10, 7, 4, 1, 14, 11, 8, 5, 2, 15, 12, 9, 6, 3}
}

func rc() []uint {
	return []uint{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1B, 0x36}
}

func g(word uint32, round uint) uint32 {
	v_0 := uint(word >> 24)
	v_1 := uint(word >> 16 & 0xFF)
	v_2 := uint(word >> 8 & 0xFF)
	v_3 := uint(word & 0xFF)

	v_0 = sbox()[v_0]
	v_1 = sbox()[v_1]
	v_2 = sbox()[v_2]
	v_3 = sbox()[v_3]

	// V1 xor rc[round]
	v_1 ^= rc()[round-1]

	// (v1,v2,v3,v0) isntead of (v0,v1,v2,v3)
	return uint32(v_1<<24 ^ v_2<<16 ^ v_3<<8 ^ v_0)
}

func mul123(input, factor uint) uint {
	if factor == 1 {
		return input
	} else if factor == 2 {
		c := input << 1
		if ((input >> 7) & 1) == 1 {
			c ^= 0x11b
		}
		return c
	} else if factor == 3 {
		return mul123(input, 2) ^ input
	}
	panic("Unreachable Code")
}

func mc(input [4]uint) uint32 {
	res := make([]uint, 4)
	res[0] = mul123(input[0], 2) ^ mul123(input[1], 3) ^ mul123(input[2], 1) ^ mul123(input[3], 1)
	res[1] = mul123(input[0], 1) ^ mul123(input[1], 2) ^ mul123(input[2], 3) ^ mul123(input[3], 1)
	res[2] = mul123(input[0], 1) ^ mul123(input[1], 1) ^ mul123(input[2], 2) ^ mul123(input[3], 3)
	res[3] = mul123(input[0], 3) ^ mul123(input[1], 1) ^ mul123(input[2], 1) ^ mul123(input[3], 2)
	return uint32(res[0]<<24 ^ res[1]<<16 ^ res[2]<<8 ^ res[3])
}

func invert_key_schedule_aes_128(round_key_10 []uint) []uint {
	roundkeys := make([][4]uint32, 11)
	for i := 0; i < 4; i++ {
		roundkeys[10][i] = uint32(round_key_10[i*4]<<24 ^ round_key_10[i*4+1]<<16 ^ round_key_10[i*4+2]<<8 ^ round_key_10[i*4+3])
	}

	for i := 10; i > 0; i-- {
		// compute in reverse order
		for j := 3; j > 0; j-- {
			roundkeys[i-1][j] = roundkeys[i][j] ^ roundkeys[i][j-1]
		}

		roundkeys[i-1][0] = roundkeys[i][0] ^ g(roundkeys[i-1][3], uint(i))
	}

	main_key := make([]uint, 16)
	for i := 0; i < 4; i++ {
		main_key[i*4] = uint(roundkeys[0][i] >> 24)
		main_key[i*4+1] = uint((roundkeys[0][i] << 8) >> 24)
		main_key[i*4+2] = uint((roundkeys[0][i] << 16) >> 24)
		main_key[i*4+3] = uint((roundkeys[0][i] << 24) >> 24)
	}

	return main_key
}

/*	 
* ATTACK HELPER
*/

// Returns true when "e" is in (a >> "shifts" for a in "s")
func cont_shifted(s [255]uint32, e uint32, shifts uint32) bool {
	for _, a := range s {
		if (a >> shifts) == e {
			return true
		}
	}
	return false
}

// Returns true when value is in one of the delta slices
// the shift value is needed to convert from uint32 to the currently attacked keybyte (1-4)
func is_in_deltas(value uint32, shifts uint32, deltas [4][255]uint32) bool {
	if cont_shifted(deltas[0], value, shifts) || cont_shifted(deltas[1], value, shifts) {
		return true
	}
	if cont_shifted(deltas[2], value, shifts) || cont_shifted(deltas[3], value, shifts) {
		return true
	}
	return false
}

/*	 
* ATTACK
*/

// Finally!
func decrypt_flag(key []uint, enc_flag []uint) []byte {
	finalKeyB := sliceToBytes(key)
	cipher, _ := aes.NewCipher(finalKeyB)
	res := make([]byte, 48)

	for i := 0; i < 3; i++ {
		flagB := sliceToBytes(enc_flag[i*16:(i+1)*16])
		pt := make([]byte, 16)
		cipher.Decrypt(pt, flagB)
		copy(res[i*16:(i+1)*16], pt)
	}
	
	return res
}

// Tests whether the found main key (!) is correct
func testEncryption(key []uint, pt []uint, ct []uint) bool {
	finalKeyB := sliceToBytes(key)
	ptB := sliceToBytes(pt)
	ctB := sliceToBytes(ct)
	res := make([]byte, 16)

	cipher, _ := aes.NewCipher(finalKeyB)

	cipher.Encrypt(res, ptB)

	// compare every byte
	correct := true
	for i := range res {
		if res[i] != ctB[i] {
			correct = false
		}
	}
	return correct
}

// Attacks all four keybytes of the column at once (used only after candidates are already ruled out!)
func attack_all_bytes(key_space map[uint32]uint32, ct, faulty []uint, deltas [4][255]uint32, column int) {
	key_space_copy := copyMap(key_space)
	for guess := range key_space_copy {
		k := []uint32{guess >> 24, (guess >> 16) & 0xFF, (guess >> 8) & 0xFF, guess & 0xFF}
		diff := make([]uint, 4)

		for i := 0; i < 4; i++ {
			correct_in := sub_bytes_input(ct, k[i], inv_shift_rows()[column*4+i])
			fault_in := sub_bytes_input(faulty, k[i], inv_shift_rows()[column*4+i])
			diff[i] = correct_in ^ fault_in
		}

		finalDiff := uint32(diff[0]<<24 ^ diff[1]<<16 ^ diff[2]<<8 ^ diff[3])

		if is_in_deltas(finalDiff, 0, deltas) {
			continue
		} else {
			// kick out candidate
			delete(key_space, uint32(guess))
		}
	}
}

// does what the name says
func sub_bytes_input(data []uint, guess uint32, position uint) uint {
	return inv_sbox()[data[position]^uint(guess)]
}

func read_input() (plaintext [][]uint, ciphertext [][]uint, faulty [][]uint, enc_flag []uint) {
	pt, err := os.ReadFile("ciphertext.json")
	check(err)
	json.Unmarshal(pt, &ciphertext)

	ct, err := os.ReadFile("plaintext.json")
	check(err)
	json.Unmarshal(ct, &plaintext)
	
	ft, err := os.ReadFile("faulty_ciphertext.json")
	check(err)
	json.Unmarshal(ft, &faulty)

	fl, err := os.ReadFile("secret_ingredient.json")
	check(err)
	json.Unmarshal(fl, &enc_flag) 

	return plaintext, ciphertext, faulty, enc_flag
}

// this function is targeted at the fault model used: for every column exactly one byte is faulty
// since it is totally random which one, we compute all possible deltas and apply mix columns
// note that we use 32bit uint here and later shift the values accordingly
func precompute_mc() (mcDeltas [4][255]uint32) {
	for guess := uint(1); guess < 256; guess++ {
		deltas := [4][4]uint{{guess, 0, 0, 0}, {0, guess, 0, 0}, {0, 0, guess, 0}, {0, 0, 0, guess}}

		for i := 0; i < 4; i++ {
			mcDeltas[i][guess-1] = mc(deltas[i])
		}

	}
	return mcDeltas
}

func main() {
	start := time.Now()
	pt, ct, faulty, enc_flag := read_input()
	elapsed := time.Since(start)
	fmt.Printf("Input read took %s\n\n", elapsed)

	// Precomputations
	mix_column_deltas := precompute_mc()
	round_key_10 := []uint{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Atack each column after another to recover Roundkey 10
	fmt.Println("Starting fault attack on roundkey 10")
	for column := 0; column < 4; column++ {
		start = time.Now()

		// Fill Key Space for first two bytes
		key_space_k12 := make(map[uint32]uint32, 65536)
		for i := uint32(0); i < 65536; i++ {
			key_space_k12[i] = uint32(i)
		}

		// First attack byte 1 and byte 2 of column
		key_space_k12_copy := copyMap(key_space_k12)
		for guess := range key_space_k12_copy {
			k := []uint32{guess >> 8, guess & 0xFF}
			diff := make([]uint, 2)

			// compute the difference between the correct input of the 10th round subbytes
			// and the faulty one
			for i := 0; i < 2; i++ {
				correct_in := sub_bytes_input(ct[0], k[i], inv_shift_rows()[column*4+i])
				faulty_in := sub_bytes_input(faulty[0], k[i], inv_shift_rows()[column*4+i])
				diff[i] = correct_in ^ faulty_in
			}

			finalDiff := uint32(diff[0]<<8 ^ diff[1])

			// then check whether the difference is in one of the deltas
			if is_in_deltas(finalDiff, 16, mix_column_deltas) {
				continue
			} else {
				// kick out candidate if not
				delete(key_space_k12, uint32(guess))
			}
		}

		// extend attack to byte 3
		key_space_k123 := make(map[uint32]uint32, 256*len(key_space_k12))
		for k3 := uint(0); k3 < 256; k3++ {
			for k12 := range key_space_k12 {
				id := uint32(key_space_k12[uint32(k12)]<<8 ^ uint32(k3))
				key_space_k123[id] = id
			}
		}

		key_space_k123_copy := copyMap(key_space_k123)
		for guess := range key_space_k123_copy {
			k := []uint32{guess >> 16, (guess >> 8) & 0xFF, guess & 0xFF}
			diff := make([]uint, 3)

			for i := 0; i < 3; i++ {
				correct_in := sub_bytes_input(ct[0], k[i], inv_shift_rows()[column*4+i])
				faulty_in := sub_bytes_input(faulty[0], k[i], inv_shift_rows()[column*4+i])
				diff[i] = correct_in ^ faulty_in
			}

			finalDiff := uint32(diff[0]<<16 ^ diff[1]<<8 ^ diff[2])

			if is_in_deltas(finalDiff, 8, mix_column_deltas) {
				continue
			} else {
				delete(key_space_k123, uint32(guess))
			}
		}

		// finally extend attack to byte 4
		key_space := make(map[uint32]uint32, 256*len(key_space_k123))
		for k4 := uint(0); k4 < 256; k4++ {
			for k123 := range key_space_k123 {
				id := uint32(key_space_k123[uint32(k123)]<<8 ^ uint32(k4))
				key_space[id] = id
			}
		}

		currentTrace := 0
		for len(key_space) > 1 && currentTrace < len(ct) {
			attack_all_bytes(key_space, ct[currentTrace], faulty[currentTrace], mix_column_deltas, column)

			currentTrace++
		}

		if len(key_space) != 1 {
			fmt.Println("ABORT: Did not find new bytes for current column.")
			os.Exit(0)
		}

		// yeah i know, a lot of shifting going on here
		for k := range key_space {
			round_key_10[column*4] = uint(k >> 24)
			round_key_10[column*4+1] = uint((k << 8) >> 24)
			round_key_10[column*4+2] = uint((k << 16) >> 24)
			round_key_10[column*4+3] = uint((k << 24) >> 24)
		}

		elapsed = time.Since(start)
		fmt.Printf("Found new bytes: %d, needed %d pairs, took %s\n", round_key_10, currentTrace, elapsed)
	}

	// dont forget to apply shift rows!
	round_key := make([]uint, 16)
	for i := 0; i < 16; i++ {
		round_key[i] = round_key_10[shift_rows()[i]]
	}
	fmt.Printf("\nRoundkey 10: %d", round_key)

	// use roundkey 10 to compute main key
	main_key := invert_key_schedule_aes_128(round_key)
	fmt.Printf("\nMain Key: %d\n", main_key)

	// test main key
	fmt.Println("\nVeryfing key...")
	for i := 0; i < 5; i++ {
		test := testEncryption(main_key, pt[i], ct[i])
		if !test {
			fmt.Printf("Encryption test %d...Fail\n", i+1)
			os.Exit(0)
		}
	}
	fmt.Println("Passed all tests")

	// decrypt Flag
	fmt.Println("\nDecrypt flag...")
	res := decrypt_flag(main_key, enc_flag)
	fmt.Printf("%s\n", res)

	// Profit
}