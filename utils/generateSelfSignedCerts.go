package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"math/big"
	"net"
	"os"
	"time"
)

var certsPath string

func init() {
	flag.StringVar(&certsPath, "path", "certs/", "Path where certificates will be stored")
}

func main() {

	flag.Parse()

	hosts := []string{"localhost", "127.0.0.1", "0.0.0.0", "local-sepp.example.com", "local-sbi.example.com"}

	caCert, caKey := generateCert("CA", nil, nil, true, hosts)

	err := writeCertAndKey(certsPath+"ca.crt", caCert, certsPath+"ca.key", caKey)
	if err != nil {
		panic(err)
	}

	n32Cert, n32Key := generateCert("N32 Server", caCert, caKey, false, hosts)
	err = writeCertAndKey(certsPath+"n32Server.crt", n32Cert, certsPath+"n32Server.key", n32Key)
	if err != nil {
		panic(err)
	}

	sbiCert, sbiKey := generateCert("SBI Server", caCert, caKey, false, hosts)
	err = writeCertAndKey(certsPath+"sbiServer.crt", sbiCert, certsPath+"sbiServer.key", sbiKey)
	if err != nil {
		panic(err)
	}

	clientCert, clientKey := generateCert("Client", caCert, caKey, false, hosts)
	err = writeCertAndKey(certsPath+"client.crt", clientCert, certsPath+"client.key", clientKey)
	if err != nil {
		panic(err)
	}

}

func generateCert(commonName string, caCert *x509.Certificate, caKey *rsa.PrivateKey, isCA bool, hosts []string) (*x509.Certificate, *rsa.PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	var certDER []byte
	if isCA {
		certDER, err = x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	} else {
		certDER, err = x509.CreateCertificate(rand.Reader, &template, caCert, &key.PublicKey, caKey)
	}
	if err != nil {
		panic(err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		panic(err)
	}

	return cert, key
}

func writeCertAndKey(certPath string, cert *x509.Certificate, keyPath string, key *rsa.PrivateKey) error {
	certOut, err := os.Create(certPath)
	if err != nil {
		panic(err)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err != nil {
		return err
	}

	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err != nil {
		return err
	}

	keyOut.Close()

	return nil
}
