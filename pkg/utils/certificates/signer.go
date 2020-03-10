/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
