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

package crypto

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

type scryptHasher struct {
	CostFactor      int
	BlockSize       int
	Parallelization int
	KeyLength       int
}

func NewScryptHasher(costFactor, blockSize, parallelization, keyLength int) PasswordHasher {
	return scryptHasher{
		CostFactor:      costFactor,
		BlockSize:       blockSize,
		Parallelization: parallelization,
		KeyLength:       keyLength,
	}
}

func (s scryptHasher) HashPassword(password string, salt string) (string, error) {
	if len(salt) == 0 {
		return "", errors.New("salt should not be empty")
	}

	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}

	saltBytes := bytes.NewBufferString(salt).Bytes()

	derivedKey, err := scrypt.Key([]byte(password), saltBytes, s.CostFactor, s.BlockSize, s.Parallelization, s.KeyLength)
	if err != nil {
		return "", err
	}

	base64Salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	base64Hash := base64.RawStdEncoding.EncodeToString(derivedKey)

	encodedHash := fmt.Sprintf("$scrypt$costFactor=%d$blockSize=%d$parallelization=%d$keyLength=%d$%s$%s", s.CostFactor, s.BlockSize, s.Parallelization, s.KeyLength, base64Salt, base64Hash)
	return encodedHash, nil
}

func (s scryptHasher) ComparePassword(password string, hashedPassword string, salt string) (bool, error) {
	costFactor, blockSize, parallelization, keyLength, decodedSalt, storedHash, err := s.decodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	computedHash, err := scrypt.Key([]byte(password), decodedSalt, costFactor, blockSize, parallelization, keyLength)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare(storedHash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (s scryptHasher) decodeHash(encodedHash string) (costFactor, blockSize, parallelization, keyLength int, salt []byte, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 8 {
		err = errors.New("invalid hash format")
		return
	}

	_, err = fmt.Sscanf(parts[2], "costFactor=%d", &costFactor)
	if err != nil {
		return
	}

	_, err = fmt.Sscanf(parts[3], "blockSize=%d", &blockSize)
	if err != nil {
		return
	}

	_, err = fmt.Sscanf(parts[4], "parallelization=%d", &parallelization)
	if err != nil {
		return
	}

	_, err = fmt.Sscanf(parts[5], "keyLength=%d", &keyLength)
	if err != nil {
		return
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[6])
	if err != nil {
		return
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[7])
	if err != nil {
		return
	}

	return
}
