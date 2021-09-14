package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	contents, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Printf("Failed to read file %s,err:%s", *filename, err.Error())
		os.Exit(1)
	}
	fmt.Println(string(contents))
}
