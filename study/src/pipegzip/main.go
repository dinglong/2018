package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

/**
 * 使用pipe做中间过渡，将write和reader连接起来
 */
func main() {
	in, err := os.Open("in.dat")
	if err != nil {
		log.Fatalf("open file error, %v\n", err)
	}

	done := make(chan bool, 1)

	r, w := io.Pipe()
	go func(ir io.Reader) {
		out, err := os.Create("out.dat.gz")
		if err != nil {
			log.Fatalf("create out file error, %v\n", err)
		}

		size, err := io.Copy(out, ir)
		if err != nil {
			log.Fatalf("copy out file error, %v\n", err)
		}
		log.Printf("copy out size %v\n", size)

		done <- true
	}(r)

	gw := gzip.NewWriter(w)
	size, err := io.Copy(gw, in)
	if err != nil {
		log.Fatalf("copy in file error, %v\n", err)
	}
	gw.Close()
	log.Printf("copy in size %v\n", size)

	w.Close()
	<-done
	r.Close()
}
