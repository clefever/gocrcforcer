package main

import (
	"fmt"
	"io"
	"math/bits"
	"os"
)

/*---- Main application ----*/

func main() {
	fmt.Println("Usage:", os.Args[0], "FileName ByteOffset NewCrc32Value")
}

/*---- Utilities ----*/

// Generator polynomial. Do not modify, because there are many dependencies
const polynomial uint64 = 0x104C11DB7

func getCrc32(raf *os.File) uint32 {
	raf.Seek(0, 0)
	var crc uint32 = 0xFFFFFFFF
	buffer := make([]byte, 32*1024)
	for {
		n, err := raf.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			return ^crc
		}
		for _, b := range buffer[:n] {
			for i := uint32(0); i < 8; i++ {
				crc ^= uint32(b>>i) << 31
				crc = (crc << 1) ^ ((crc >> 31) * uint32(polynomial&0xFFFFFFFF))
			}
		}
	}
}

func reverseBits(x uint32) uint32 {
	var y uint32
	for i := uint32(0); i < 32; i++ {
		y |= ((x >> i) & 1) << (31 - i)
	}
	return y
}

/*---- Polynomial arithmetic ----*/

// Returns polynomial x multiplied by polynomial y modulo the generator polynomial.
func multiplyMod(x uint64, y uint64) uint64 {
	// Russian peasant multiplication algorithm
	var z uint64
	for y != 0 {
		z ^= x * (y & 1)
		y >>= 1
		x <<= 1
		if (x>>32)&1 != 0 {
			x ^= polynomial
		}
	}
	return z
}

// Returns polynomial x to the power of natural number y modulo the generator polynomial.
func powMod(x uint64, y uint64) uint64 {
	// Exponentiation by squaring
	var z uint64 = 1
	for y != 0 {
		if y&1 != 0 {
			z = multiplyMod(z, x)
		}
		x = multiplyMod(x, x)
		y >>= 1
	}
	return z
}

// Computes polynomial x divided by polynomial y, returning the quotient and remainder.
func divideAndRemainder(x uint64, y uint64) (uint64, uint64) {
	if y == 0 {
		panic("Division by zero")
	}
	if x == 0 {
		return 0, 0
	}

	ydeg := getDegree(y)
	var z uint64
	for i := getDegree(x) - ydeg; i >= 0; i-- {
		if ((x >> uint64(i+ydeg)) & 1) != 0 {
			x ^= y << uint64(i)
			z |= uint64(1) << uint64(i)
		}
	}
	return z, x
}

// Returns the reciprocal of polynomial x with respect to the generator polynomial.
func reciprocalMod(x uint64) uint64 {
	// Based on a simplification of the extended Euclidean algorithm
	y := x
	x = polynomial
	var a uint64
	var b uint64 = 1
	for y != 0 {
		q, r := divideAndRemainder(x, y)
		c := a ^ multiplyMod(q, b)
		x = y
		y = r
		a = b
		b = c
	}
	if x != 1 {
		panic("Reciprocal does not exist")
	}
	return a
}

func getDegree(x uint64) int32 {
	return int32(bits.Len64(x) - 1)
}
