package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := flag.String("f", "", "Specify the file name")
	flag.Parse()

	if filename == nil || len(*filename) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Printf("Failed to read file %s,err:%s", *filename, err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 64)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Failed to read,err:%s", err.Error())
			break
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buf[:n]))
	}
}
