package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	cert := flag.String("cert", "cert.pem", "cert file path")
	private := flag.String("private", "private.pem", "private key file path")
	host := flag.String("host", "localhost", "cert hostname")
	expiresIn := flag.Duration("expires-in", 365*24*time.Hour, "cert expires in")
	org := flag.String("org", "go-advanced", "cert organization name")

	flag.Parse()

	if err := genTLS(*cert, *private, *host, *expiresIn, *org); err != nil {
		log.Fatalf("generate tls error: %s", err.Error())
	}

	fmt.Printf("OK:\n")
	fmt.Printf("  certificate: %s\n", *cert)
	fmt.Printf("  private key: %s\n", *private)
}

func genTLS(certPath, privatePath, host string, expiresIn time.Duration, organization string) error {
	cert := &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   host,
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(expiresIn),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("serial number error: %w", err)
	}

	cert.SerialNumber = serialNumber

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("private key error: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("create certificate err: %w", err)
	}

	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return fmt.Errorf("encode certificate error: %w", err)
	}

	var privatePEM bytes.Buffer
	if err = pem.Encode(&privatePEM, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return fmt.Errorf("encode private error: %w", err)
	}

	if err := os.WriteFile(certPath, certPEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("write certificate error: %w", err)
	}

	if err := os.WriteFile(privatePath, privatePEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("write private error: %w", err)
	}

	return nil
}
