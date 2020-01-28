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
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io"
	"os"
	"time"
)

// Signer provides a client certificate signer
type Signer interface {
	// GenerateClient generates a client certificate for us
	GenerateClient(string, string, time.Duration) ([]byte, []byte, error)
	// GenerateServer generates a server certificate for us
	GenerateServer([]string, time.Duration) ([]byte, []byte, error)
}

type signImpl struct {
	ca  *x509.Certificate
	key *rsa.PrivateKey
}

// GenerateClient is used to generate a client certificate
func (s *signImpl) GenerateClient(cn, ou string, duration time.Duration) ([]byte, []byte, error) {
	return CreateCertificate(s.ca, s.key, x509.Certificate{
		Subject: pkix.Name{
			CommonName:         cn,
			OrganizationalUnit: []string{ou},
		},
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:              time.Now().Add(duration),
		NotBefore:             time.Now().Add(-30 * time.Second),
	})
}

// GenerateServer is used to generate a server certificate
func (s *signImpl) GenerateServer(hosts []string, duration time.Duration) ([]byte, []byte, error) {
	return CreateCertificate(s.ca, s.key, x509.Certificate{
		Subject: pkix.Name{
			CommonName:         hosts[0],
			OrganizationalUnit: []string{"Appvia Kore"},
		},
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              hosts,
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:              time.Now().Add(duration),
		NotBefore:             time.Now().Add(-30 * time.Second),
	})
}

// NewSignerFromFiles creates and returns a signer from files
func NewSignerFromFiles(authority, key string) (Signer, error) {
	cafile, err := os.Open(authority)
	if err != nil {
		return nil, err
	}
	cakey, err := os.Open(key)
	if err != nil {
		return nil, err
	}

	return NewSigner(cafile, cakey)
}

// NewSigner creates and returns a new signer
func NewSigner(authority, key io.Reader) (Signer, error) {
	// @step: load the certificates
	ca, err := LoadCertificateAuthority(authority)
	if err != nil {
		return nil, err
	}
	privatekey, err := LoadPrivateKey(key)
	if err != nil {
		return nil, err
	}

	return &signImpl{ca: ca, key: privatekey}, nil
}
