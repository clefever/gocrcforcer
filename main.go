package main

import (
	"fmt"
	"math/bits"
	"os"
)

func main() {
	fmt.Println("Usage:", os.Args[0], "FileName ByteOffset NewCrc32Value")
}

func getDegree(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}
