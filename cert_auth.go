package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"os"
	"strings"
	"time"
)

var commonNameCA = "Devcert Certificate Authority (CA)"

// CA represents a certificate authority cert and private key.
type CA struct {
	Valid   bool
	Cert    *x509.Certificate
	PrivKey *rsa.PrivateKey
}



func buildCAPaths() (crtPath, keyPath string, err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		err = fmt.Errorf("Building CA paths failed: %w", err)
		return
	}

	var crtB, keyB strings.Builder
	crtB.WriteString(devcertDir)
	crtB.WriteString("devcert_ca.crt")

	keyB.WriteString(devcertDir)
	keyB.WriteString("devcert_ca.key")

	crtPath = crtB.String()
	keyPath = keyB.String()

	return
}

// loadCA will load the certificate authority data from the files.
func loadCA() (ca *CA, err error) {

	certPath, keyPath, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("Loading CA failed: %w", err)
		return
	}

	caExist, err := caFilesExist()
	if err != nil {
		err = fmt.Errorf("Loading CA failed: %w", err)
		return
	}

	ca = &CA{
		Valid: false,
	}

	// Files doesn't exist
	if caExist == false {
		return
	}

	crtBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return ca, nil
	}

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return ca, nil
	}

	crtBlock, _ := pem.Decode(crtBytes)
	ca.Cert, err = x509.ParseCertificate(crtBlock.Bytes)
	if err != nil {
		return ca, nil
	}

	keyBlock, _ := pem.Decode(keyBytes)
	ca.PrivKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return ca, nil
	}

	// Check if the CA has expired
	now := time.Now()
	if ca.Cert.NotBefore.After(now) == true || ca.Cert.NotAfter.Before(now) == true {
		return ca, nil
	}

	// Certificate is valid
	ca.Valid = true

	return ca, nil
}

// caFilesExist check if the certificate authority files exist.
func caFilesExist() (exist bool, err error) {
	certPath, keyPath, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("Checking CA files exist failed: %w", err)
		return
	}

	// Check cert file
	_, err = os.Stat(certPath)
	certNotExist := errors.Is(err, fs.ErrNotExist)

	if err != nil && certNotExist == false {
		err = fmt.Errorf("Checking CA files exist failed: %w", err)
		return
	}

	// Check key file
	_, err = os.Stat(keyPath)
	keyNotExist := errors.Is(err, fs.ErrNotExist)

	if err != nil && keyNotExist == false {
		err = fmt.Errorf("Checking CA files exist failed: %w", err)
		return
	}

	err = nil
	if certNotExist == false && keyNotExist == false {
		exist = true
	}

	return
}

func createCA() (err error) {
	fmt.Printf("Creating certificate authority (CA) files...\n")

	ca, err := loadCA()
	if err != nil {
		err = fmt.Errorf("Creating CA failed: %w", err)
		return
	}

	// Certificate is valid, nothing to do.
	if ca.Valid == true {
		fmt.Printf("Certificate authority (CA) already created.\n")
		return
	}

	// Certificate is invalid and the CA files exists, remove and re-generate it.

	exist, err := caFilesExist()
	if err != nil {
		err = fmt.Errorf("Creating CA failed: %w", err)
		return
	}

	// Remove files if exist
	if exist == true {
		err = removeCAFiles()
		if err != nil {
			err = fmt.Errorf("Creating CA failed: %w", err)
			return
		}
	}

	// Create CA
	err = generateCA()
	if err != nil {
		fmt.Errorf("Creating CA failed: %w", err)
		return
	}

	return
}

// generateCA creates the certificate authority files.
func generateCA() (err error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(mrand.Int63()),
		Subject: pkix.Name{
			CommonName: commonNameCA,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = fmt.Errorf("Generating CA failed: %w", err)
		return
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		err = fmt.Errorf("Generating CA failed: %w", err)
		return
	}

	// Create the files on the file system
	crtPath, keyPath, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("Generating CA failed: %w", err)
		return
	}

	caFile, err := os.Create(crtPath)
	if err != nil {
		err = fmt.Errorf("Generating CA failed: %w", err)
		return
	}

	defer caFile.Close()

	caPrivKeyFile, err := os.Create(keyPath)
	if err != nil {
		err = fmt.Errorf("Generating CA failed: %w", err)
		return
	}

	defer caPrivKeyFile.Close()

	caPEMWriter := bufio.NewWriter(caFile)
	pem.Encode(caPEMWriter, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPEMWriter.Flush()

	caPrivKeyPEMWriter := bufio.NewWriter(caPrivKeyFile)
	pem.Encode(caPrivKeyPEMWriter, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	caPrivKeyPEMWriter.Flush()

	fmt.Printf("Certificate authority (CA) created at\n  Certificate: %s\n  Private Key: %s\n", crtPath, keyPath)

	return
}

// removeCAFiles removes the CA crt and key files from the devcert folder.
func removeCAFiles() (err error) {
	crtPath, keyPath, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("Removing CA files failed: %w", err)
		return
	}

	err = os.Remove(crtPath)
	if err != nil {
		err = fmt.Errorf("Removing CA files failed: %w", err)
		return
	}

	err = os.Remove(keyPath)
	if err != nil {
		err = fmt.Errorf("Removing CA files failed: %w", err)
		return
	}

	return
}
