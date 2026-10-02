package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Private Key
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateFile, err := os.Create("utils/private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()

	pem.Encode(privateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	})

	// Public Key
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	publicFile, err := os.Create("utils/public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()

	pem.Encode(publicFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})

	fmt.Println("RSA key berhasil dibuat")
}
