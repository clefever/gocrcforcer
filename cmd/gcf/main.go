package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/clefever/gocrcforcer"
)

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
	err = gocrcforcer.ModifyFileCrc32(os.Args[1], offset, uint32(newcrc), true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
