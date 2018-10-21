package main

import (
	"errors"
	"fmt"
	"io"
	"math/bits"
	"os"
	"strconv"
)

/*---- Main application ----*/

func main() {
	// Handle arguments
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "FileName ByteOffset NewCrc32Value")
		os.Exit(1)
	}

	// Parse and check file offset argument
	offset, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Invalid byte offset")
		os.Exit(1)
	}

	// Parse and check new CRC argument
	if len(os.Args[3]) != 8 {
		fmt.Fprintln(os.Stderr, "Error: Invalid new CRC-32 value")
		os.Exit(1)
	}
	newcrc, err := strconv.ParseUint(os.Args[3], 16, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Invalid new CRC-32 value")
		os.Exit(1)
	}

	// Process the file
	err = ModifyFileCrc32(os.Args[1], offset, bits.Reverse32(uint32(newcrc)), true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

/*---- Main function ----*/

// ModifyFileCrc32 - Public library function.
func ModifyFileCrc32(path string, offset int64, newcrc uint32, printstatus bool) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	filelen := fi.Size()
	if offset+4 > filelen {
		return errors.New("Byte offset plus 4 exceeds file length")
	}

	// Read entire file and calculate original CRC-32 value
	crc, err := getCrc32(f)
	if err != nil {
		return err
	}
	if printstatus {
		fmt.Printf("Original CRC-32: %X\n", bits.Reverse32(crc))
	}

	// Compute the change to make
	rm, err := reciprocalMod(powMod(2, uint64((filelen-offset)*8)))
	if err != nil {
		return err
	}
	delta := uint32(multiplyMod(rm, uint64(crc^newcrc)))

	// Patch 4 bytes in the file
	if _, err := f.Seek(offset, 0); err != nil {
		return err
	}
	bytes4 := make([]byte, 4)
	if _, err := f.Read(bytes4); err != nil {
		return err
	}
	for i := range bytes4 {
		bytes4[i] ^= uint8(bits.Reverse32(delta) >> uint32(i*8))
	}
	if _, err := f.Seek(offset, 0); err != nil {
		return err
	}
	if _, err := f.Write(bytes4); err != nil {
		return err
	}
	if printstatus {
		fmt.Println("Computed and wrote patch")
	}

	// Recheck entire file
	if chkcrc, err := getCrc32(f); err != nil || chkcrc != newcrc {
		return errors.New("Failed to update CRC-32 to desired value")
	}
	if printstatus {
		fmt.Println("New CRC-32 successfully verified")
	}
	return nil
}

/*---- Utilities ----*/

// Generator polynomial. Do not modify, because there are many dependencies
const polynomial uint64 = 0x104C11DB7

func getCrc32(f *os.File) (uint32, error) {
	f.Seek(0, 0)
	var crc uint32 = 0xFFFFFFFF
	buffer := make([]byte, 32*1024)
	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			return ^crc, nil
		}
		for _, b := range buffer[:n] {
			for i := uint32(0); i < 8; i++ {
				crc ^= uint32(b>>i) << 31
				crc = (crc << 1) ^ ((crc >> 31) * uint32(polynomial&0xFFFFFFFF))
			}
		}
	}
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
func reciprocalMod(x uint64) (uint64, error) {
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
		return 0, errors.New("Reciprocal does not exist")
	}
	return a, nil
}

func getDegree(x uint64) int32 {
	return int32(bits.Len64(x) - 1)
}
