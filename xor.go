// Package xor implements simple XOR based ciphers.
package xor

import (
	"fmt"
	"io"
)

// Create a new xor stream cipher, key bytes will be read from keyStream.
func NewReadWriter(rw io.ReadWriter, keyStream io.Reader) io.ReadWriter {
	return readWriter{
		writer: NewWriter(rw, keyStream).(writer),
		reader: NewReader(rw, keyStream).(reader),
	}
}

// NewReader returns a new Reader that XOR'es bytes from r with bytes from keyStream.
func NewReader(r, keyStream io.Reader) io.Reader {
	if r == nil {
		panic("NewReader: reader is nil")
	}
	if keyStream == nil {
		panic("NewReader: keyStream is nil")
	}

	return reader{
		Reader:    r,
		keyStream: keyStream,
	}
}

func NewWriter(w io.Writer, keyStream io.Reader) io.Writer {
	if w == nil {
		panic("NewWriter: writer is nil")
	}
	if keyStream == nil {
		panic("NewWriter: keyStream is nil")
	}

	return writer{
		Writer:    w,
		keyStream: keyStream,
	}
}

// readWriter Wraps a Writer and Reader.
type readWriter struct {
	writer
	reader
}

// writer wraps another io.Writer but ciphers data with XOR.
type writer struct {
	io.Writer           // underlying Writer
	keyStream io.Reader // where to read key bytes from
}

func (w writer) Write(plaintext []byte) (n int, err error) {
	// Read the amount of key bytes needed.
	key := make([]byte, len(plaintext))
	_, err = io.ReadFull(w.keyStream, key)
	if err != nil {
		return 0, fmt.Errorf("Read from keyStream: %v", err)
	}

	// Buffer to hold ciphertext
	ciphertext := make([]byte, len(plaintext))
	n = BufXOR(ciphertext, plaintext, key)
	if n < len(ciphertext) {
		return 0, fmt.Errorf("XOR Write: wrote %d; should have wrote %d",
			n, len(ciphertext))
	}

	// Now write the ciphertext to the underlying Writer
	return w.Writer.Write(ciphertext)
}

type reader struct {
	io.Reader           // underlying Reader
	keyStream io.Reader // where to read key bytes from
}

func (r reader) Read(dst []byte) (n int, err error) {
	// read len(dst) into ciphertext
	ciphertext := make([]byte, len(dst))
	_, err = io.ReadFull(r.Reader, ciphertext)
	if err != nil {
		return 0, fmt.Errorf("Read failed: %v", err)
	}

	// get enough key bytes from the keyStream
	key := make([]byte, len(dst))
	_, err = io.ReadFull(r.keyStream, key)
	if err != nil {
		return 0, fmt.Errorf("Reading keyStream: %v", err)
	}

	// XOR them together
	return BufXOR(dst, ciphertext, key), nil
}

// XOR two byte arrays together
func XOR(a, b []byte) []byte {
	size := len(a)
	if len(b) > size {
		size = len(b)
	}

	buf := make([]byte, size)
	n := BufXOR(buf, a, b)
	return buf[:n]
}

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
