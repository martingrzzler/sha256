package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/bits"
	"os"
)

const CHUNK_SIZE = 64
const CHUNK_SIZE_BITS = CHUNK_SIZE * 8
const POW2_32 = 4294967296

//  L + 1 + K + 64 % 512 = 0
func NumPaddZero(L int) int {
	lenInBits := L * 8

	m := lenInBits + 1 + 64
	return CHUNK_SIZE_BITS - m%CHUNK_SIZE_BITS
}

func PaddMessage(data string) []uint8 {
	b := []uint8(data)
	dataLen := len(data)
	zerosLen := ((NumPaddZero(dataLen) + 1) / 8) - 1
	firstByteAfterMessage := uint8(0b10000000)

	b = append(b, firstByteAfterMessage)
	for i := 0; i < zerosLen; i++ {
		b = append(b, 0b00000000)
	}

	lenAsBytes := make([]uint8, 8)
	binary.BigEndian.PutUint64(lenAsBytes, uint64(dataLen*8))
	b = append(b, lenAsBytes...)

	return b
}

func Chunks(data []uint8) [][]uint8 {
	chunks := make([][]uint8, 0)

	for i := 0; i < len(data); i += CHUNK_SIZE {
		chunks = append(chunks, data[i:i+CHUNK_SIZE])
	}

	return chunks
}

func MessageSchedule(chunk []uint8) []uint32 {
	messageSchedule := make([]uint32, 0)
	for i := 0; i < len(chunk); i += 4 {
		messageSchedule = append(messageSchedule, binary.BigEndian.Uint32(chunk[i:i+4]))
	}

	for j := 16; j < 64; j++ {
		messageSchedule = append(messageSchedule, 0)
	}

	return messageSchedule
}

func main() {
	message := os.Args[1]

	initial := []uint32{
		0x6a09e667,
		0xbb67ae85,
		0x3c6ef372,
		0xa54ff53a,
		0x510e527f,
		0x9b05688c,
		0x1f83d9ab,
		0x5be0cd19,
	}
	k := []uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

	padded := PaddMessage(message)
	chunks := Chunks(padded)

	for _, chunk := range chunks {
		ms := MessageSchedule(chunk)

		for i := 16; i < 64; i++ {
			s0 := bits.RotateLeft32(ms[i-15], -7) ^ bits.RotateLeft32(ms[i-15], -18) ^ (ms[i-15] >> 3)
			s1 := bits.RotateLeft32(ms[i-2], -17) ^ bits.RotateLeft32(ms[i-2], -19) ^ (ms[i-2] >> 10)
			ms[i] = ms[i-16] + s0 + ms[i-7] + s1
		}

		a := initial[0]
		b := initial[1]
		c := initial[2]
		d := initial[3]
		e := initial[4]
		f := initial[5]
		g := initial[6]
		h := initial[7]

		for j := 0; j < 64; j++ {
			S1 := bits.RotateLeft32(e, -6) ^ bits.RotateLeft32(e, -11) ^ bits.RotateLeft32(e, -25)
			ch := (e & f) ^ ((^e) & g)
			temp1 := h + S1 + ch + k[j] + ms[j]
			S0 := bits.RotateLeft32(a, -2) ^ bits.RotateLeft32(a, -13) ^ bits.RotateLeft32(a, -22)
			maj := (a & b) ^ (a & c) ^ (b & c)
			temp2 := S0 + maj

			h = g
			g = f
			f = e
			e = d + temp1
			d = c
			c = b
			b = a
			a = temp1 + temp2
		}
		initial[0] = initial[0] + a
		initial[1] = initial[1] + b
		initial[2] = initial[2] + c
		initial[3] = initial[3] + d
		initial[4] = initial[4] + e
		initial[5] = initial[5] + f
		initial[6] = initial[6] + g
		initial[7] = initial[7] + h
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, initial)

	digest := hex.EncodeToString(buf.Bytes())
	fmt.Println(digest)
}
