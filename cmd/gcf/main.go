package main

import (
	"log"
	"os"
	"strconv"

	"github.com/clefever/gocrcforcer"
)

func main() {
	// Handle arguments
	if len(os.Args) != 4 {
		log.Fatalln("Usage:", os.Args[0], "FileName ByteOffset NewCrc32Value")
	}

	// Parse and check file offset argument
	offset, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalln("Error: Invalid byte offset")
	}

	// Parse and check new CRC argument
	if len(os.Args[3]) != 8 {
		log.Fatalln("Error: Invalid new CRC-32 value")
	}
	newcrc, err := strconv.ParseUint(os.Args[3], 16, 32)
	if err != nil {
		log.Fatalln("Error: Invalid new CRC-32 value")
	}

	// Process the file
	err = gocrcforcer.ModifyFileCrc32(os.Args[1], offset, uint32(newcrc), true)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
