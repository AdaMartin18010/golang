// Package main demonstrates Go 1.26 crypto/hpke package usage
//
// HPKE (Hybrid Public Key Encryption) is defined in RFC 9180.
// Go 1.26 includes crypto/hpke as a standard library package.
//
// This example demonstrates the API surface. In production,
// use proper key management and ephemeral keys for each operation.
package main

import (
	"crypto/ecdh"
	"crypto/hpke"
	"fmt"
)

func main() {
	fmt.Println("Go 1.26 crypto/hpke Demo (RFC 9180)")
	fmt.Println("====================================")

	// Select HPKE suite components
	kem := hpke.DHKEM(ecdh.X25519())
	kdf := hpke.HKDFSHA256()
	aead := hpke.AES128GCM()

	fmt.Println("\nHPKE Suite:")
	fmt.Printf("  KEM : DHKEM(X25519, HKDF-SHA256) (id=%d)\n", kem.ID())
	fmt.Printf("  KDF : %v\n", kdf)
	fmt.Printf("  AEAD: %v\n", aead)

	// Generate receiver key pair using the KEM
	hpkePriv, err := kem.GenerateKey()
	if err != nil {
		panic(err)
	}
	hpkePub := hpkePriv.PublicKey()

	// Sender encrypts a message to the receiver's public key
	info := []byte("application context info")
	plaintext := []byte("Hello, HPKE in Go 1.26!")

	ciphertext, err := hpke.Seal(hpkePub, kdf, aead, info, plaintext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nPlaintext : %s\n", plaintext)
	fmt.Printf("Ciphertext: %x...\n", ciphertext[:16])

	// Receiver decrypts with their private key
	decrypted, err := hpke.Open(hpkePriv, kdf, aead, info, ciphertext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted : %s\n", decrypted)

	fmt.Println("\ncrypto/hpke standard library APIs (Go 1.26):")
	fmt.Println("  hpke.Seal(pk, kdf, aead, info, plaintext) -> ciphertext")
	fmt.Println("  hpke.Open(k, kdf, aead, info, ciphertext) -> plaintext")
	fmt.Println("\nSupported KEM:")
	fmt.Println("  hpke.DHKEM(ecdh.X25519())")
	fmt.Println("  hpke.DHKEM(ecdh.P256())")
	fmt.Println("  hpke.MLKEM768()")
	fmt.Println("  hpke.MLKEM1024()")
	fmt.Println("  hpke.MLKEM768X25519()  (hybrid post-quantum)")
	fmt.Println("\nSupported AEAD:")
	fmt.Println("  hpke.AES128GCM()")
	fmt.Println("  hpke.AES256GCM()")
	fmt.Println("  hpke.ChaCha20Poly1305()")
	fmt.Println("\nSupported KDF:")
	fmt.Println("  hpke.HKDFSHA256()")
	fmt.Println("  hpke.HKDFSHA384()")
	fmt.Println("  hpke.HKDFSHA512()")
}
