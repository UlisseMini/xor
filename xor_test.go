package xor

import (
	"bytes"
	"testing"
)

// Test the XOR proporties
func TestXORProporties(t *testing.T) {
	t.Run("plaintext XOR with ciphertext", func(t *testing.T) {
		// plaintext and key must be equal length
		plaintext := []byte("Plain text!")
		key := []byte("secret key!")

		ciphertext := XOR(plaintext, key)

		out := XOR(plaintext, ciphertext)
		if !bytes.Equal(out, key) {
			t.Fatalf("wanted %q; got %q", key, out)
		}
	})

	t.Run("key XOR with ciphertext", func(t *testing.T) {
		// plaintext and key must be equal length
		plaintext := []byte("Plain text!")
		key := []byte("secret key!")

		ciphertext := XOR(plaintext, key)

		out := XOR(key, ciphertext)
		if !bytes.Equal(out, plaintext) {
			t.Fatalf("wanted %q; got %q", plaintext, out)
		}
	})
}

// lifted off crypto/cipher/XOR.go
func TestXOR(t *testing.T) {
	for alignP := 0; alignP < 2; alignP++ {
		for alignQ := 0; alignQ < 2; alignQ++ {
			for alignD := 0; alignD < 2; alignD++ {
				p := make([]byte, 1024)[alignP:]
				q := make([]byte, 1024)[alignQ:]

				d1 := XOR(p, q)
				d2 := XOR(p, q)
				if !bytes.Equal(d1, d2) {
					t.Error("not equal")
				}
			}
		}
	}
}

// test that it works with a key larger then the data.
func TestXORBigKey(t *testing.T) {
	for alignP := 0; alignP < 2; alignP++ {
		for alignQ := 0; alignQ < 2; alignQ++ {
			for alignD := 0; alignD < 2; alignD++ {
				big := make([]byte, 4096)[alignP:]
				q := make([]byte, 1024)[alignQ:]

				d1 := XOR(q, big)
				d2 := XOR(q, big)
				if !bytes.Equal(d1, d2) {
					t.Error("not equal")
				}
			}
		}
	}
}
