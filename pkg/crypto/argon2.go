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
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type argon2IDHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewArgon2IDHasher(memory uint32, iterations uint32, parallelism uint8, saltLength uint32, keyLength uint32) PasswordHasher {
	return &argon2IDHasher{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		saltLength:  saltLength,
		keyLength:   keyLength,
	}

}

func (a argon2IDHasher) HashPassword(password string, salt string) (string, error) {
	if len(salt) == 0 {
		return "", errors.New("salt should not be empty")
	}

	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}

	hash := argon2.IDKey([]byte(password), []byte(salt), a.iterations, a.memory, a.parallelism, a.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString([]byte(salt))
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.iterations, a.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (a argon2IDHasher) ComparePassword(password string, decodeHash string, salt string) (bool, error) {
	memory, iterations, parallelism, decodedSalt, hash, err := a.decodeHash(decodeHash)
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), decodedSalt, iterations, memory, parallelism, a.keyLength)

	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}
	return false, nil

}

func (a argon2IDHasher) decodeHash(encodedHash string) (memory uint32, iterations uint32, parallelism uint8, salt []byte, hash []byte, err error) {
	var version int

	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		err = errors.New("invalid hash format")
		return
	}

	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return
	}
	if version != argon2.Version {
		err = errors.New("invalid argon2 version")
		return
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return
	}

	return
}
