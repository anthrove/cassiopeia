/*
 * Copyright (C) 2025 Anthrove
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"time"
)

func GenerateRSAKey(bitSize int) (*rsa.PrivateKey, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, bitSize)

	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	return key, keyPEM, nil
}

func GenerateECDSAKey(curve elliptic.Curve) (*ecdsa.PrivateKey, []byte, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PEM format.
	privateKeyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return key, privateKeyPem, nil
}

func GenerateCertificate(privateKey any, publicKey any, expireAt time.Time, signatureAlgo x509.SignatureAlgorithm) ([]byte, error) {
	tml := x509.Certificate{
		// you can add any attr that you need
		NotBefore: time.Now(),
		NotAfter:  expireAt,
		// TODO you have to generate a different serial number each execution
		SerialNumber: big.NewInt(123123),
		Subject: pkix.Name{
			CommonName:   "New DisplayName",
			Organization: []string{"New Org."},
		},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    signatureAlgo,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, publicKey, privateKey)
	if err != nil {
		log.Fatal("Certificate cannot be created.", err.Error())
	}

	// Generate a pem block with the certificate
	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	return certPem, nil
}

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	if len(pub) == 0 {
		return nil, errors.New("empty public key")
	}
	block, _ := pem.Decode(pub)
	if block == nil {
		return nil, errors.New("public key error")
	}
	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}
	return key, nil
}
