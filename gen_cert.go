package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// generateCertificate creates the certificate using the domain names.
func generateCertificate(domains []string) (err error) {
	fmt.Printf("Generating certificate...\n")

	isValid := true
	var invalidDomains []string

	for _, domain := range domains {
		valid := validateDomain(domain)
		if valid == false {
			isValid = false
			invalidDomains = append(invalidDomains, domain)
		}
	}

	if isValid == false {
		fmt.Printf("The folowing domain names are invalid:\n")

		for _, d := range invalidDomains {
			fmt.Printf("  - %s\n", d)
		}

		err = fmt.Errorf("")
		return
	}

	commonName := []string{"Devcert Certificate - "}
	commonName = append(commonName, domains...)

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(mrand.Int63()),
		Subject: pkix.Name{
			CommonName: strings.Join(commonName, ""),
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 11, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     domains,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	ca, err := loadCA()
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &certPrivKey.PublicKey, ca.PrivKey)
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	devcertDir, err := buildDevcertDir()
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	var crtB, keyB strings.Builder
	crtB.WriteString(devcertDir)
	crtB.WriteString("devcert_")
	crtB.WriteString(domains[0])
	crtB.WriteString("_multi.crt")

	keyB.WriteString(devcertDir)
	keyB.WriteString("devcert_")
	keyB.WriteString(domains[0])
	keyB.WriteString("_multi.key")

	certFile, err := os.Create(crtB.String())
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	defer certFile.Close()

	certKeyFile, err := os.Create(keyB.String())
	if err != nil {
		err = fmt.Errorf("Generating certificate failed: %w", err)
		return
	}

	defer certKeyFile.Close()

	certPEMWriter := bufio.NewWriter(certFile)
	pem.Encode(certPEMWriter, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPEMWriter.Flush()

	certKeyPEMWriter := bufio.NewWriter(certKeyFile)
	pem.Encode(certKeyPEMWriter, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	certKeyPEMWriter.Flush()

	fmt.Printf("Generated at:\n  Certificate: %s\n  Private Key: %s\n\n", crtB.String(), keyB.String())

	fmt.Printf("Valid for:\n")
	for i, domain := range domains {
		fmt.Printf("  %d. %s\n", (i + 1), domain)
	}

	return
}

func validateDomain(domain string) bool {
	regexDomain := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,255})$`)
	return regexDomain.MatchString(domain)
}
