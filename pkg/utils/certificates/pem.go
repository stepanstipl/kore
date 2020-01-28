/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"time"
)

// CreateClientCertificate used to generate a client certificate
func CreateClientCertificate(cert *x509.Certificate, signer *rsa.PrivateKey, subject string, duration time.Duration) ([]byte, []byte, error) {
	return CreateCertificate(cert, signer, x509.Certificate{
		Subject: pkix.Name{
			CommonName: subject,
		},
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:              time.Now().Add(duration),
		NotBefore:             time.Now().Add(-30 * time.Second),
	})
}

// CreateServerCertificate is used to create a server certificate
func CreateServerCertificate(cert *x509.Certificate, signer *rsa.PrivateKey, hosts []string, duration time.Duration) ([]byte, []byte, error) {
	return CreateCertificate(cert, signer, x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Appvia Kore"},
		},
		BasicConstraintsValid: true,
		DNSNames:              hosts,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:              time.Now().Add(duration),
		NotBefore:             time.Now().Add(-30 * time.Second),
	})
}

// CreateCertificate is used to sign a certificate request for us
func CreateCertificate(cert *x509.Certificate, key *rsa.PrivateKey, template x509.Certificate) ([]byte, []byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	template.SerialNumber = serialNumber

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, cert.PublicKey, key)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	buf := &bytes.Buffer{}

	if err := pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return []byte{}, []byte{}, err
	}

	priv, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	bufp := &bytes.Buffer{}
	if err := pem.Encode(bufp, &pem.Block{Type: "PRIVATE KEY", Bytes: priv}); err != nil {
		return []byte{}, []byte{}, err
	}

	return buf.Bytes(), bufp.Bytes(), nil
}

// LoadCertificateAuthority loads the CA pem file from disk
func LoadCertificateAuthority(reader io.Reader) (*x509.Certificate, error) {
	// @step: read in the context
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// @step: decode the content
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("attempting to decode the certificate")
	}

	// @step: parse the certificate
	return x509.ParseCertificate(block.Bytes)
}

// LoadPrivateKey is used to load a private key
func LoadPrivateKey(reader io.Reader) (*rsa.PrivateKey, error) {
	// @step: read in the context
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// @step: decode the content
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("attempting to decode the certificate")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
