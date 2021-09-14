package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	max_threads_number = 100
	default_chunk_size = 5 * 1024 * 1024
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
		log.Printf("Failed open file %s,err:%s", *filename, err.Error())
		return
	}
	defer f.Close()

	// buffer
	r := bufio.NewReader(f)
	// define sync.pool
	bytePool := sync.Pool{New: func() interface{} {
		return make([]byte, default_chunk_size)
	}}

	stat, err := f.Stat()
	if err != nil {
		log.Printf("Faield get file %s stat,err:%s", *filename, err.Error())
		return
	}
	filesize := stat.Size()
	log.Printf("file size is %d\n", filesize)

	// set thread number
	thNums := max_threads_number
	if (filesize / default_chunk_size) < max_threads_number {
		thNums = int(filesize / default_chunk_size)
	}
	tokenQueue := make(chan struct{}, thNums)
	errChan := make(chan error, 1)

	for i := 0; i < thNums; i++ {
		tokenQueue <- struct{}{}
	}

	errs := make([]error, 0)
	errWrapper := make(chan error, 1)
	go func() {
		defer func() {
			//pass err to out func
			if len(errs) > 0 {
				errWrapper <- fmt.Errorf("%v", errs)
			}
			close(errWrapper)
		}()

		for deleteErr := range errChan {
			errs = append(errs, deleteErr)
		}
	}()

	var wg sync.WaitGroup
	usedSize := 0
	for {
		<-tokenQueue
		buf := bytePool.Get().([]byte)
		n, err := r.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				fmt.Println("end of file")
			} else if err != nil {
				log.Printf("Failed read file %s,err:%s", *filename, err.Error())
			}
			break
		}
		usedSize = len(buf) + usedSize
		p := float64(usedSize) / float64(filesize)
		log.Printf("Percent:%f\n", p)

		wg.Add(1)
		// use thread process data
		go func() {
			defer func() {
				tokenQueue <- struct{}{}
				wg.Done()
			}()
			// process data
			if err := process(); err != nil {
				errChan <- err
			}
		}()
	}

	// wait all goroutine end
	wg.Wait()
	// close chan
	close(errChan)

	if err := <-errChan; err != nil {
		log.Printf("process file %s err:%s", *filename, err.Error())
		return
	}
	t2 := time.Now()
	log.Printf("end time:%s", t2.String())
}

func process() error {
	return nil
}
