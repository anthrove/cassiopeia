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
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"strings"
)

type pbkdf2Hasher struct {
	Digest     func() hash.Hash
	Iterations int
	KeyLen     int
}

func NewPBKDF2Hasher(digest func() hash.Hash, iterations int, keyLen int) PasswordHasher {
	return pbkdf2Hasher{
		Digest:     digest,
		Iterations: iterations,
		KeyLen:     keyLen,
	}
}

func (p pbkdf2Hasher) HashPassword(password string, salt string) (string, error) {
	if len(salt) == 0 {
		return "", errors.New("salt should not be empty")
	}

	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}

	saltBuffer := bytes.NewBufferString(salt).Bytes()

	df := pbkdf2.Key([]byte(password), saltBuffer, p.Iterations, p.KeyLen, p.Digest)

	b64Salt := base64.RawStdEncoding.EncodeToString(saltBuffer)
	b64Hash := base64.RawStdEncoding.EncodeToString(df)

	encodedHash := fmt.Sprintf("$pbkdf2$iter=%d$keylen=%d$%s$%s", p.Iterations, p.KeyLen, b64Salt, b64Hash)
	return encodedHash, nil
}

func (p pbkdf2Hasher) ComparePassword(password string, hashedPassword string, salt string) (bool, error) {
	iterations, keyLen, decodedSalt, passwordHash, err := p.decodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	computedHash := pbkdf2.Key([]byte(password), decodedSalt, iterations, keyLen, p.Digest)

	if subtle.ConstantTimeCompare(passwordHash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (p pbkdf2Hasher) decodeHash(encodedHash string) (iterations int, keyLen int, salt []byte, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		err = errors.New("invalid hash format")
		return
	}

	_, err = fmt.Sscanf(vals[2], "iter=%d", &iterations)
	if err != nil {
		return
	}

	_, err = fmt.Sscanf(vals[3], "keylen=%d", &keyLen)
	if err != nil {
		return
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return
	}

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return
	}

	return
}
