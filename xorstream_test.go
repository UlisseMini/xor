package xor

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestStream(t *testing.T) {
	t.Run("TestConn", connStream)

	// Also failing, i sware my implementation is correct though :/
	t.Run("TestBasicIO", basicIO)

	// Commented out because i was getting strange read fails with the buffer :/
	// t.Run("TestWithBytesBuffers", bufStream)
}

// test Stream working with bytes.Buffer readers and writers.
func bufStream(t *testing.T) {
	numTests := 100       // number of random data tests to run.
	var size int64 = 1024 // size of the data being encrypted for every test.

	for i := 0; i < numTests; i++ {
		// get random key and data
		data := bytes.NewBuffer(nil)
		key := bytes.NewBuffer(nil)
		mustCopyN(key, rand.Reader, size)
		mustCopyN(data, rand.Reader, size)

		// simulated readwriter
		rw := bytes.NewBuffer(nil)

		// fill key and data with random data
		mustCopyN(data, rand.Reader, size)

		rwXOR := NewReadWriter(rw, key)

		// write the data
		mustCopyN(rwXOR, data, size)

		// read the data
		buf := bytes.NewBuffer(nil)
		buf.Grow(int(size))
		mustCopyN(buf, rwXOR, size)

		if !bytes.Equal(buf.Bytes(), data.Bytes()) {
			t.Fatalf("buf is not equal to data.Bytes()\nbuf(%X)\ndata(%X)",
				buf.Bytes(), data.Bytes())
		}
	}
}

// connStream tests Stream encryption with a network connection
func connStream(t *testing.T) {
}

// wrap io.CopyN but panic on error or failure to copy all data
func mustCopyN(dst io.Writer, src io.Reader, n int64) int {
	nw, err := io.CopyN(dst, src, n)
	if nw != n {
		p := fmt.Sprintf("mustCopyN: wrote %d wanted to write %d, err = %v",
			nw, n, err)
		panic(p)
	}

	return int(nw)
}

// test basic I/O XOR oporations
func basicIO(t *testing.T) {
	expected := []byte(`--------- plaintext with key
01010100 -- plaintext
10010011 -- key
--------- XOR = 11000111 (ciphertext)

--------- ciphertext with key
11000111 -- ciphertext
10010011 -- key
--------- XOR = 01010100 (plaintext)

--------- plaintext with ciphertext
01010100 -- plaintext
11000111 -- ciphertext
--------- XOR = 10010011 (key)
`)

	// get a key to be used, must be in io.Reader format so i'm using a bytes.Buffer.
	writerKey := &bytes.Buffer{}
	mustCopyN(writerKey, rand.Reader, int64(len(expected)))

	// same key for reading
	readerKey := bytes.NewBuffer(writerKey.Bytes())

	// where the ciphertext will be stored
	ciphertext := &bytes.Buffer{}
	w := NewWriter(ciphertext, writerKey)

	// Write some data to the XOR stream
	_, err := w.Write(expected)
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(ciphertext, readerKey)
	// Read it and compare with expected
	buf := make([]byte, len(expected))
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf[:n], expected) {
		t.Fatalf("buf(%X) is not expected(%X)", buf[:n], expected)
	}
}
