package main

import (
	"bufio"
	"flag"
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
		log.Printf("Failed to open file %s,err:%s", *filename, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Printf("Failed to get file %s stat,err:%s", *filename, err.Error())
		os.Exit(1)
	}
	log.Printf("The file %s size is %d", *filename, stat.Size())

	//刚才还在想为什么它每次读取没有设置 offset，因为第一读取的 offset=0，下次读取就是从上次读取之后的位置再往后读取，第一次读取了三个字节，
	//那么下次读取就是从第四个字节开始读取了
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadBytes('\n')
		log.Print(string(line))

		if err != nil {
			if err == io.EOF {
				log.Println("end of file")
				break
			} else {
				log.Printf("read file err:%s", err.Error())
				break
			}
		}
	}
}
