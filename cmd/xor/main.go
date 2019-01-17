// Package main provides a commandline fronted for the XOR package.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[*] ")
}

func main() {
	// by default take one input from stdin and output to stdout
	out := os.Stdout
	bufSize := 32 * 1024

	inFileName := flag.String("i", "",
		"File to read input from.")
	flag.Parse()

	if *inFileName == "" {
		// TODO: print flag help instead
		log.Print("Must supply -i")
		return
	}

	infile, err := os.Open(*inFileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer infile.Close()

	if err := StreamXOR(os.Stdin, infile, out, bufSize); err != nil {
		log.Println(err)
		return
	}
}

// XOR two data streams into an output stream,
// if any of it fails it will return an error.
//
// If buf is nil, one will be automatically allocated.
func StreamXOR(in1, in2 io.Reader, out io.Writer, size int) (err error) {
	buf1 := make([]byte, 32*1024)   // in1 buffer
	buf2 := make([]byte, 32*1024)   // in2 buffer
	outbuf := make([]byte, 32*1024) // output buffer

	for {
		// Fill input buffers
		n1r, err1 := in1.Read(buf1)
		n2r, err2 := in2.Read(buf2)
		if (err1 != nil && err1 != io.EOF) || (err2 != nil && err2 != io.EOF) {
			return fmt.Errorf("read failed: err1 = %v; err2 = %v",
				err1, err2)
		}

		// XOR into outbuf
		n := BufXOR(outbuf, buf1[:n1r], buf2[:n2r])
		_, err := out.Write(outbuf[:n])
		if err != nil {
			return fmt.Errorf("write failed: %v", err)
		}
	}

	return err
}

// "A little copying is better then a little dependency"
// - Rob Pike
func BufXOR(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}
